package main

import (
	"fmt"
	"github.com/gramarr/sonarr"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func (e *Env) HandleTVShowDetails(m *tb.Message) {
	e.CM.StartConversation(NewTVShowDetailsConversation(e), m)
}

func NewTVShowDetailsConversation(e *Env) *TVShowDetailsConversation {
	return &TVShowDetailsConversation{env: e}
}

type TVShowDetailsConversation struct {
	currentStep Handler
	env         *Env
	tvShow      sonarr.TVShow
}

func (c *TVShowDetailsConversation) Run(m *tb.Message) {
	c.currentStep = c.showTVShowDetails(m)
}

func (c *TVShowDetailsConversation) Name() string {
	return "tvdetails"
}

func (c *TVShowDetailsConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *TVShowDetailsConversation) showTVShowDetails(m *tb.Message) Handler {
	tvShowId, err := strconv.Atoi(m.Payload)
	if err != nil {
		return nil
	}

	c.tvShow, err = c.env.Sonarr.GetTVShow(tvShowId)
	if err != nil {
		return nil
	}

	c.tvShow.RemotePoster = c.env.Sonarr.GetPosterURL(c.tvShow)
	if c.tvShow.RemotePoster == "" {
		tvDBTVShow, _ := c.env.Sonarr.SearchTVShow(c.tvShow.TvdbID)
		c.tvShow.RemotePoster = c.env.Sonarr.GetPosterURL(tvDBTVShow)
	}

	if c.tvShow.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(c.tvShow.RemotePoster)}
		_, _ = c.env.Bot.Send(m.Sender, photo)
	}

	var msg []string
	msg = append(msg, fmt.Sprintf("*%s (%d)*", EscapeMarkdown(c.tvShow.Title), c.tvShow.Year))
	msg = append(msg, c.tvShow.Overview)
	msg = append(msg, "")
	msg = append(msg, fmt.Sprintf("*Previous Airing:* %s", FormatDate(c.tvShow.PreviousAiring)))
	msg = append(msg, fmt.Sprintf("*Next Airing:* %s", FormatDate(c.tvShow.NextAiring)))
	msg = append(msg, fmt.Sprintf("*Folder:* %s", GetRootFolderFromPath(c.tvShow.Path)))
	requesterList := c.env.Sonarr.GetRequesterList(c.tvShow)
	msg = append(msg, fmt.Sprintf("*Requested by:* %s", strings.Join(requesterList, ", ")))

	Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))

	var username = GetUserName(m)

	var options []string
	user, exists := c.env.Users.User(m.Sender.ID)
	if exists {
		if user.IsAdmin() {
			if len(c.tvShow.Tags) > 0 {
				options = append(options, "Remove requester")
			}
			options = append(options, "Add requester")
			options = append(options, "Delete tv show")
		} else {
			for _, t := range c.tvShow.Tags {
				tag, _ := c.env.Sonarr.GetTagById(t)
				if tag.Label == username {
					options = append(options, "Remove yourself from the requester list")
					break
				}
			}
			if len(options) == 0 {
				options = append(options, "Add yourself to the requester list")
			}
		}
	}

	options = append(options, "Back to tv show list")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Choose an option for this tv show", options)

	return func(m *tb.Message) {
		switch m.Text {
		case "Add yourself to the requester list":
			c.addRequester(m, username)
		case "Remove yourself from the requester list":
			c.removeRequester(m, username)
		case "Remove requester":
			c.currentStep = c.askRemoveRequester(m)
		case "Add requester":
			c.currentStep = c.askAddRequester(m)
		case "Back to tv show list":
			c.env.HandleListTVShows(m)
			c.env.CM.StopConversation(c)
		case "Delete tv show":
			err = c.env.Sonarr.DeleteTVShow(c.tvShow.ID)
			if err == nil {
				Send(c.env.Bot, m.Sender, fmt.Sprintf("TVShow '%s (%d)' has been deleted.", EscapeMarkdown(c.tvShow.Title), c.tvShow.Year))
			} else {
				SendError(c.env.Bot, m.Sender, "Could not delete tv show. Please try again.")
			}
		default:
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			SendKeyboardList(c.env.Bot, m.Sender, "Choose an option for this tv show", options)
		}
	}
}

func (c *TVShowDetailsConversation) askRemoveRequester(m *tb.Message) Handler {
	var options []string
	for _, t := range c.tvShow.Tags {
		tag, err := c.env.Sonarr.GetTagById(t)
		if err == nil {
			options = append(options, tag.Label)
		}
	}
	options = append(options, "Back to tv show details")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Remove user from requester list", options)

	return func(m *tb.Message) {
		if m.Text != "Back to tv show details" {
			c.removeRequester(m, m.Text)
		}
		m.Payload = strconv.Itoa(c.tvShow.ID)
		c.currentStep = c.showTVShowDetails(m)
	}
}

func (c *TVShowDetailsConversation) askAddRequester(m *tb.Message) Handler {
	tags, err := c.env.Sonarr.GetTags()
	var options []string
	if err == nil {
		for _, t := range tags {
			options = append(options, t.Label)
		}
	} else {
		SendError(c.env.Bot, m.Sender, "Could not retrieve tag list")
	}
	options = append(options, "Back to tv show details")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Add user from requester list", options)

	return func(m *tb.Message) {
		if m.Text != "Back to tv show tv show details" {
			c.addRequester(m, m.Text)
		}
		m.Payload = strconv.Itoa(c.tvShow.ID)
		c.currentStep = c.showTVShowDetails(m)
	}
}

func (c *TVShowDetailsConversation) removeRequester(m *tb.Message, requester string) {
	var err error
	c.tvShow, err = c.env.Sonarr.RemoveRequester(c.tvShow, requester)
	if err == nil {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("%s has been removed from the requester list.", requester))
		if len(c.tvShow.Tags) == 0 {
			SendAdmin(c.env.Bot, c.env.Users.Admins(),
				fmt.Sprintf("'%s' was the last requester of the tvShow '%s (%d)'. Send /deletetv\\_%d to delete it from disk.",
					requester, EscapeMarkdown(c.tvShow.Title), c.tvShow.Year, c.tvShow.ID))
		}
	} else {
		Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
	}
}

func (c *TVShowDetailsConversation) addRequester(m *tb.Message, requester string) {
	var err error
	c.tvShow, err = c.env.Sonarr.AddRequester(c.tvShow, requester)
	if err == nil {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("%s has been added to the requester list.", requester))
	} else {
		Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
	}
}
