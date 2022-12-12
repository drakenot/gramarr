package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/internal/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) HandleAuth(m *tb.Message) {
	var msg []string
	pass := m.Payload
	user, exists := mb.users.User(m.Sender.ID)

	// Empty Password?
	if pass == "" {
		mb.tele.Send(m.Sender, "Usage: `/auth [password]`")
		return
	}

	// Is User Already Admin?
	if exists && user.IsAdmin() {
		// Notify User
		msg = append(msg, "You're already authorized.")
		msg = append(msg, "Type /start to begin.")
		mb.tele.Send(m.Sender, strings.Join(msg, "\n"))
		return
	}

	// Check if pass is Admin Password
	if pass == mb.config.Bot.AdminPassword {
		if exists {
			user.Access = users.UAAdmin
			mb.users.Update(user)
		} else {
			newUser := users.User{
				ID:        m.Sender.ID,
				FirstName: m.Sender.FirstName,
				LastName:  m.Sender.LastName,
				Username:  m.Sender.Username,
				Access:    users.UAAdmin,
			}
			mb.users.Create(newUser)
		}

		// Notify User
		msg = append(msg, "You have been authorized as an *admin*.")
		msg = append(msg, "Type /start to begin.")
		mb.tele.Send(m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted admin access.", util.DisplayName(m.Sender))
		mb.tele.SendAdmin(mb.users.Admins(), adminMsg)

		return
	}

	// Check if pass is User Password
	if pass == mb.config.Bot.Password {
		if exists {
			// Notify User
			msg = append(msg, "You're already authorized.")
			msg = append(msg, "Type /start to begin.")
			mb.tele.Send(m.Sender, strings.Join(msg, "\n"))
			return
		}
		newUser := users.User{
			ID:        m.Sender.ID,
			Username:  m.Sender.Username,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			Access:    users.UAMember,
		}
		mb.users.Create(newUser)

		// Notify User
		msg = append(msg, "You have been authorized.")
		msg = append(msg, "Type /start to begin.")
		mb.tele.Send(m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted acccess.", util.DisplayName(m.Sender))
		mb.tele.SendAdmin(mb.users.Admins(), adminMsg)
		return
	}

	// Notify User
	mb.tele.SendError(m.Sender, "Your password is invalid.")

	// Notify Admin
	adminMsg := "%s made an invalid auth request with password: %s"
	adminMsg = fmt.Sprintf(adminMsg, util.DisplayName(m.Sender), util.EscapeMarkdown(m.Payload))
	mb.tele.SendAdmin(mb.users.Admins(), adminMsg)
}
