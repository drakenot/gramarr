package main

import (
	"fmt"
	"github.com/gramarr/sonarr"
	tb "gopkg.in/tucnak/telebot.v2"
	"sort"
	"strconv"
	"strings"
)

func (e *Env) HandleListTVShows(m *tb.Message) {
	e.CM.StartConversation(NewListTVShowsConversation(e), m)
}

func NewListTVShowsConversation(e *Env) *ListTVShowsConversation {
	return &ListTVShowsConversation{env: e}
}

type ListTVShowsConversation struct {
	currentStep            Handler
	tvShowQuery            string
	tvShowResults          []sonarr.TVShow
	folderResults          []sonarr.Folder
	selectedTVShow         *sonarr.TVShow
	selectedQualityProfile *sonarr.Profile
	selectedFolder         *sonarr.Folder
	env                    *Env
}

func (c *ListTVShowsConversation) Run(m *tb.Message) {
	c.currentStep = c.AskRequester(m)
}

func (c *ListTVShowsConversation) Name() string {
	return "listTVShows"
}

func (c *ListTVShowsConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *ListTVShowsConversation) AskRequester(m *tb.Message) Handler {

	requesterList, err := c.env.Sonarr.GetTags()

	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get requester list.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(requesterList) == 0 {
		SendError(c.env.Bot, m.Sender, "No requester found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	var options []string
	options = append(options, "All")
	for _, requester := range requesterList {
		options = append(options, fmt.Sprintf("%s", strings.TrimSpace(requester.Label)))
	}

	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which tv shows would you like to list?", options)

	return func(m *tb.Message) {
		// Set the selected folder
		for _, opt := range options {
			if m.Text == opt {
				if m.Text == "All" {
					c.tvShowResults, _ = c.env.Sonarr.GetTVShows()
				} else {
					c.tvShowResults, err = c.env.Sonarr.GetTVShowsByRequester(strings.TrimSpace(m.Text))
				}
			}
		}

		if err != nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskRequester(m)
			return
		}

		c.currentStep = c.AskTVShow(m)
	}
}

func (c *ListTVShowsConversation) AskTVShow(m *tb.Message) Handler {

	sort.Slice(c.tvShowResults, func(i, j int) bool {
		return c.tvShowResults[i].Title < c.tvShowResults[j].Title
	})

	var pending = []string{"*Requested TV Shows:*"}
	var options []string
	for _, tvShow := range c.tvShowResults {
		options = append(options, EscapeMarkdown(tvShow.Title))
		pending = append(pending, fmt.Sprintf("- %s", EscapeMarkdown(tvShow.Title)))
	}
	if len(pending) > 1 {
		Send(c.env.Bot, m.Sender, strings.Join(pending, "\n"))
	}

	if len(options) > 0 {
		options = append(options, "Back to requester selection")
		options = append(options, "/cancel")
		SendKeyboardList(c.env.Bot, m.Sender, "Select a tv show for more details", options)
	} else {
		Send(c.env.Bot, m.Sender, "No tv shows found")
		c.currentStep = c.AskRequester(m)
	}

	return func(m *tb.Message) {
		if m.Text == "Back to requester selection" {
			c.currentStep = c.AskRequester(m)
			return
		}

		// Set the selected tvShow
		for i, opt := range options {
			if m.Text == opt {
				c.selectedTVShow = &c.tvShowResults[i]
				break
			}
		}

		// Not a valid tvShow selection
		if c.selectedTVShow == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			return
		}

		m.Payload = strconv.Itoa(c.selectedTVShow.ID)
		c.env.HandleTVShowDetails(m)
		c.env.CM.StopConversation(c)
	}

}
