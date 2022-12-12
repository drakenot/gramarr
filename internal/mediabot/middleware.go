package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/internal/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) RequireAuth(access users.UserAccess, h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		user, _ := mb.users.User(m.Sender.ID)
		var msg []string

		// Is Revoked?
		if user.IsRevoked() {
			// Notify User
			msg = append(msg, "Your access has been revoked and you cannot reauthorize.")
			msg = append(msg, "Please reach out to the bot owner for support.")
			mb.tele.SendError(m.Sender, strings.Join(msg, "\n"))

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Revoked user %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			mb.tele.SendAdmin(mb.users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Not Member?
		isAuthorized := user.IsAdmin() || user.IsMember()
		if !isAuthorized && access != users.UANone {
			// Notify User
			mb.tele.SendError(m.Sender, "You are not authorized to use this bot.\n`/auth [password]` to authorize.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("Unauthorized user %s attempted the following command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			mb.tele.SendAdmin(mb.users.Admins(), strings.Join(msg, "\n"))
			return
		}

		// Is Non-Admin and requires Admin?
		if !user.IsAdmin() && access == users.UAAdmin {
			// Notify User
			mb.tele.SendError(m.Sender, "Only admins can use this command.")

			// Notify Admins
			msg = append(msg, fmt.Sprintf("User %s attempted the following admin command:", util.DisplayName(m.Sender)))
			msg = append(msg, fmt.Sprintf("`%s`", m.Text))
			mb.tele.SendAdmin(mb.users.Admins(), strings.Join(msg, "\n"))
			return
		}

		h(m)
	}
}

func (mb *MediaBot) RequirePrivate(h func(m *tb.Message)) func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}
		h(m)
	}
}
