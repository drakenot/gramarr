package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

func (e *Env) HandleUsers(m *tb.Message) {
	err := e.Users.Load()
	if err != nil {
		Send(e.Bot, m.Sender, "Error loading users")
		return
	}

	var msg []string

	admins := e.Users.Admins()
	if len(admins) > 0 {
		msg = append(msg, "*Admins:*")
		for i := range admins {
			if len(admins[i].Username) > 0 {
				msg = append(msg, admins[i].Username)
			} else {
				msg = append(msg, admins[i].FirstName)
			}
		}
	}

	users := e.Users.Users()
	if len(users) > 0 {
		msg = append(msg, "\n*Users:*")
		for i := range users {
			if !users[i].IsAdmin() {
				if len(users[i].Username) > 0 {
					msg = append(msg, users[i].Username)
				} else {
					msg = append(msg, users[i].FirstName)
				}
			}
		}
	}

	if len(msg) > 0 {
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
	}
}
