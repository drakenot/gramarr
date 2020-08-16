package main

import (
	"fmt"
	"github.com/gramarr/radarr"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func (e *Env) HandleDetails(m *tb.Message) {
	e.CM.StartConversation(NewDetailsConversation(e), m)
}

func NewDetailsConversation(e *Env) *DetailsConversation {
	return &DetailsConversation{env: e}
}

type DetailsConversation struct {
	currentStep Handler
	env         *Env
	movie       radarr.Movie
}

func (c *DetailsConversation) Run(m *tb.Message) {
	c.currentStep = c.showDetails(m)
}

func (c *DetailsConversation) Name() string {
	return "details"
}

func (c *DetailsConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *DetailsConversation) showDetails(m *tb.Message) Handler {
	movieId, err := strconv.Atoi(m.Payload)
	if err != nil {
		return nil
	}

	c.movie, err = c.env.Radarr.GetMovie(movieId)
	if err != nil {
		return nil
	}

	if c.movie.RemotePoster == "" {
		traktMovie, _ := c.env.Radarr.SearchMovie(c.movie.TmdbID)
		c.movie.RemotePoster = c.env.Radarr.GetPosterURL(traktMovie)
	}

	if c.movie.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(c.movie.RemotePoster)}
		_, _ = c.env.Bot.Send(m.Sender, photo)
	}

	var msg []string
	msg = append(msg, fmt.Sprintf("*%s (%d)*", EscapeMarkdown(c.movie.Title), c.movie.Year))
	msg = append(msg, c.movie.Overview)
	msg = append(msg, "")
	msg = append(msg, fmt.Sprintf("*Cinema Date:* %s", FormatDate(c.movie.InCinemas)))
	msg = append(msg, fmt.Sprintf("*BluRay Date:* %s", FormatDate(c.movie.PhysicalRelease)))
	msg = append(msg, fmt.Sprintf("*Folder:* %s", GetRootFolderFromPath(c.movie.Path)))
	if c.movie.HasFile {
		msg = append(msg, fmt.Sprintf("*Downloaded:* %s", FormatDateTime(c.movie.MovieFile.DateAdded)))
		msg = append(msg, fmt.Sprintf("*File:* %s", c.movie.MovieFile.RelativePath))
	} else {
		msg = append(msg, "*Downloaded:* No")
	}
	requesterList := c.env.Radarr.GetRequesterList(c.movie)
	msg = append(msg, fmt.Sprintf("*Requested by:* %s", strings.Join(requesterList, ", ")))

	Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))

	var username string
	if len(m.Sender.Username) > 0 {
		username = m.Sender.Username
	} else {
		username = m.Sender.FirstName
	}

	var options []string
	user, exists := c.env.Users.User(m.Sender.ID)
	if exists {
		if user.IsAdmin() {
			if len(c.movie.Tags) > 0 {
				options = append(options, "Remove requester")
			}
			options = append(options, "Add requester")
			options = append(options, "Delete Movie")
		} else {
			for _, t := range c.movie.Tags {
				tag, _ := c.env.Radarr.GetTagById(t)
				if tag.Label == strings.ToLower(username) {
					options = append(options, "Remove yourself from the requester list")
					break
				}
			}
			if len(options) == 0 {
				options = append(options, "Add yourself to the requester list")
			}
		}
	}

	options = append(options, "Back to movie list")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Choose an option for this movie", options)

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
		case "Back to movie list":
			c.env.HandleListMovies(m)
			c.env.CM.StopConversation(c)
		case "Delete Movie":
			err = c.env.Radarr.DeleteMovie(c.movie.ID)
			if err == nil {
				Send(c.env.Bot, m.Sender, fmt.Sprintf("Movie '%s (%d)' has been deleted.", EscapeMarkdown(c.movie.Title), c.movie.Year))
			} else {
				SendError(c.env.Bot, m.Sender, "Could not delete movie. Please try again.")
			}
		default:
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			SendKeyboardList(c.env.Bot, m.Sender, "Choose an option for this movie", options)
		}
	}
}

func (c *DetailsConversation) askRemoveRequester(m *tb.Message) Handler {
	var options []string
	for _, t := range c.movie.Tags {
		tag, err := c.env.Radarr.GetTagById(t)
		if err == nil {
			options = append(options, tag.Label)
		}
	}
	options = append(options, "Back to details")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Remove user from requester list", options)

	return func(m *tb.Message) {
		if m.Text != "Back to details" {
			c.removeRequester(m, m.Text)
		}
		m.Payload = strconv.Itoa(c.movie.ID)
		c.currentStep = c.showDetails(m)
	}
}

func (c *DetailsConversation) askAddRequester(m *tb.Message) Handler {
	tags, err := c.env.Radarr.GetTags()
	var options []string
	if err == nil {
		for _, t := range tags {
			options = append(options, t.Label)
		}
	} else {
		SendError(c.env.Bot, m.Sender, "Could not retrieve tag list")
	}
	options = append(options, "Back to details")
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Add user from requester list", options)

	return func(m *tb.Message) {
		if m.Text != "Back to details" {
			c.addRequester(m, m.Text)
		}
		m.Payload = strconv.Itoa(c.movie.ID)
		c.currentStep = c.showDetails(m)
	}
}

func (c *DetailsConversation) removeRequester(m *tb.Message, requester string) {
	var err error
	c.movie, err = c.env.Radarr.RemoveRequester(c.movie, requester)
	if err == nil {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("%s has been removed from the requester list.", requester))
		if len(c.movie.Tags) == 0 {
			SendAdmin(c.env.Bot, c.env.Users.Admins(),
				fmt.Sprintf("'%s' was the last requester of the movie '%s (%d)'. Send /deletemovie\\_%d to delete it from disk.",
					requester, EscapeMarkdown(c.movie.Title), c.movie.Year, c.movie.ID))
		}
	} else {
		Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
	}
}

func (c *DetailsConversation) addRequester(m *tb.Message, requester string) {
	var err error
	c.movie, err = c.env.Radarr.AddRequester(c.movie, requester)
	if err == nil {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("%s has been added to the requester list.", requester))
	} else {
		Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
	}
}
