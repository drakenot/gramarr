package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"fmt"
)

func (e *Env) RequireAuth(access UserAccess, h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		user, _ := e.Users.User(m.Sender.ID)
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			SendError(e.Bot, m.Sender, strings.Join(msg, "\n"))

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked user %s attempted the following command:", DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Not Member?
		isAuthorized := user.IsAdmin() || user.IsMember()
		if !isAuthorized && access != UANone {
			// Notify User
			SendError(e.Bot, m.Sender, "You are not authorized to use this bot.\n`/auth [password]` to authorize.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Unauthorized user %s attempted the following command:", DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == UAAdmin {
			// Notify User
			SendError(e.Bot, m.Sender, "Only admins can use this command.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			SendAdmin(e.Bot, e.Users.Admins(), strings.Join(msg, "\n"))
			return
		}

		h(m)
	}
}

func (e *Env) RequirePrivate(h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}
		h(m)
	}
}
