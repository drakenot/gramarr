package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

func (e *Env) HandleDeleteTVShow(m *tb.Message) {
	tvShowId, err := strconv.Atoi(m.Payload)
	if err != nil {
		Send(e.Bot, m.Sender, "Please enter a valid tv show ID")
		return
	}

	tvShow, err := e.Sonarr.GetTVShow(tvShowId)
	if err != nil {
		Send(e.Bot, m.Sender, fmt.Sprintf("Tv show with ID '%s' not found:", m.Payload))
		return
	}

	err = e.Sonarr.DeleteTVShow(tvShowId)
	if err != nil {
		Send(e.Bot, m.Sender, "Something went wrong. Please try again")
		return
	} else {
		Send(e.Bot, m.Sender, fmt.Sprintf("Tv show %s (%d) has been deleted", EscapeMarkdown(tvShow.Title), tvShow.Year))
	}
}
