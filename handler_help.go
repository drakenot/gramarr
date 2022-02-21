package main

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleHelp(m *tb.Message) {

	user, exists := e.Users.User(m.Sender.ID)

	var msg []string
	msg = append(msg, fmt.Sprintf("Hello, I'm %s! Use these commands to control me:", e.Bot.Me.FirstName))

	if !exists {
		msg = append(msg, "")
		msg = append(msg, "/auth [password] - authenticate with the bot")
	}

	if exists && user.IsAdmin() {
		msg = append(msg, "")
		msg = append(msg, "*Admin*")
		msg = append(msg, "/users - list all bot users")
		msg = append(msg, "/status - displays Sonarr/Radarr server status")
		msg = append(msg, "/moviedetails {id} - displays movie details")
		msg = append(msg, "/deletemovie {id} - deletes movie")
		msg = append(msg, "/tvdetails {id} - displays tv show details")
		msg = append(msg, "/deletetv {id} - deletes tv show")
	}

	if exists && (user.IsMember() || user.IsAdmin()) {
		msg = append(msg, "")
		msg = append(msg, "*Movies*")
		msg = append(msg, "/addmovie - add a movie")
		msg = append(msg, "/listmovies - displays all movies")
		msg = append(msg, "")
		msg = append(msg, "*TV Shows*")
		msg = append(msg, "/addtv - add a tv show")
		msg = append(msg, "/listtv - displays all tv shows")
		msg = append(msg, "")
		msg = append(msg, "/cancel - cancel the current operation")
	}

	Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
}
