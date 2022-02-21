package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

func (e *Env) HandleDeleteMovie(m *tb.Message) {
	movieId, err := strconv.Atoi(m.Payload)
	if err != nil {
		Send(e.Bot, m.Sender, "Please enter a valid movie ID")
		return
	}

	movie, err := e.Radarr.GetMovie(movieId)
	if err != nil {
		Send(e.Bot, m.Sender, fmt.Sprintf("Movie with ID '%s' not found:", m.Payload))
		return
	}

	err = e.Radarr.DeleteMovie(movieId)
	if err != nil {
		Send(e.Bot, m.Sender, "Something went wrong. Please try again")
		return
	} else {
		Send(e.Bot, m.Sender, fmt.Sprintf("Movie %s (%d) has been deleted", EscapeMarkdown(movie.Title), movie.Year))
	}
}
