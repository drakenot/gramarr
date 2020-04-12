package main

import (
	"fmt"
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

	movie, err := c.env.Radarr.GetMovie(movieId)
	if err != nil {
		return nil
	}

	if movie.RemotePoster == "" {
		traktMovie, _ := c.env.Radarr.SearchMovie(movie.TmdbID)
		movie.RemotePoster = c.env.Radarr.GetPosterURL(traktMovie)
	}

	if movie.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(movie.RemotePoster)}
		_, _ = c.env.Bot.Send(m.Sender, photo)
	}

	var msg []string
	msg = append(msg, fmt.Sprintf("*%s (%d)*", EscapeMarkdown(movie.Title), movie.Year))
	msg = append(msg, movie.Overview)
	msg = append(msg, "")
	msg = append(msg, fmt.Sprintf("*Cinema Date:* %s", FormatDate(movie.InCinemas)))
	msg = append(msg, fmt.Sprintf("*BluRay Date:* %s", FormatDate(movie.PhysicalRelease)))
	msg = append(msg, fmt.Sprintf("*Folder:* %s", GetRootFolderFromPath(movie.Path)))
	if movie.HasFile {
		msg = append(msg, fmt.Sprintf("*Downloaded:* %s", FormatDateTime(movie.MovieFile.DateAdded)))
		msg = append(msg, fmt.Sprintf("*File:* %s", movie.MovieFile.RelativePath))
	} else {
		msg = append(msg, "*Downloaded:* No")
	}
	msg = append(msg, fmt.Sprintf("*Requested by:* %s", strings.Join(c.env.Radarr.GetRequesterList(movie), ", ")))

	Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))

	var username string
	if len(m.Sender.Username) > 0 {
		username = m.Sender.Username
	} else {
		username = m.Sender.FirstName
	}

	var options []string
	for _, t := range movie.Tags {
		tag, _ := c.env.Radarr.GetTagById(t)
		if tag.Label == strings.ToLower(username) {
			options = append(options, "Remove yourself from the requester list")
		}
	}
	if len(options) == 0 {
		options = append(options, "Add yourself to the requester list")
	}

	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Modify requester list", options)

	return func(m *tb.Message) {
		for _, opt := range options {
			if m.Text == opt {
				if m.Text == "Add yourself to the requester list" {
					movie, err = c.env.Radarr.AddRequester(movie, username)
					if err != nil {
						Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
					} else {
						Send(c.env.Bot, m.Sender, "You have been added to the requester list.")
					}
					break
				} else if m.Text == "Remove yourself from the requester list" {
					movie, err = c.env.Radarr.RemoveRequester(movie, username)
					if err != nil {
						Send(c.env.Bot, m.Sender, "Something went wrong. Please try again.")
					} else {
						Send(c.env.Bot, m.Sender, "You have been removed from the requester list.")
						if len(movie.Tags) == 0 {
							SendAdmin(c.env.Bot, c.env.Users.Admins(),
								fmt.Sprintf("'%s' was the last requester of the movie '%s (%d)'. Send /deletemovie\\_%d to delete it from disk.",
									username, EscapeMarkdown(movie.Title), movie.Year, movie.ID))
						}
					}
					break
				}
			} else {
				SendError(c.env.Bot, m.Sender, "Invalid selection.")
			}
		}
		c.env.CM.StopConversation(c)
	}
}
