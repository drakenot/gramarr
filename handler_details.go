package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

func (e *Env) HandleDetails(m *tb.Message) {
	movieId, err := strconv.Atoi(m.Payload)
	if err != nil {
		return
	}

	movie, err := e.Radarr.GetMovie(movieId)
	if err != nil {
		return
	}

	if movie.RemotePoster == "" {
		traktMovie, _ := e.Radarr.SearchMovie(movie.TmdbID)
		movie.RemotePoster = e.Radarr.GetPosterURL(traktMovie)
	}

	if movie.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(movie.RemotePoster)}
		_, _ = e.Bot.Send(m.Sender, photo)
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
		msg = append(msg, fmt.Sprintf("*Downloaded:* %s", BoolToYesOrNo(movie.HasFile)))
	}
	msg = append(msg, fmt.Sprintf("*Requested by:* %s", strings.Join(e.Radarr.GetRequesterList(movie), ", ")))

	Send(e.Bot, m.Sender, strings.Join(msg, "\n"))

}
