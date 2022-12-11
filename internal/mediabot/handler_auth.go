package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/internal/util"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *MediaBot) HandleAuth(m *tb.Message) {
	var msg []string
	pass := m.Payload
	user, exists := b.Users.User(m.Sender.ID)

	// Empty Password?
	if pass == "" {
		util.Send(b.TClient, m.Sender, "Usage: `/auth [password]`")
		return
	}

	// Is User Already Admin?
	if exists && user.IsAdmin() {
		// Notify User
		msg = append(msg, "You're already authorized.")
		msg = append(msg, "Type /start to begin.")
		util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))
		return
	}

	// Check if pass is Admin Password
	if pass == b.Config.Bot.AdminPassword {
		if exists {
			user.Access = users.UAAdmin
			b.Users.Update(user)
		} else {
			newUser := users.User{
				ID:        m.Sender.ID,
				FirstName: m.Sender.FirstName,
				LastName:  m.Sender.LastName,
				Username:  m.Sender.Username,
				Access:    users.UAAdmin,
			}
			b.Users.Create(newUser)
		}

		// Notify User
		msg = append(msg, "You have been authorized as an *admin*.")
		msg = append(msg, "Type /start to begin.")
		util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted admin access.", util.DisplayName(m.Sender))
		util.SendAdmin(b.TClient, b.Users.Admins(), adminMsg)

		return
	}

	// Check if pass is User Password
	if pass == b.Config.Bot.Password {
		if exists {
			// Notify User
			msg = append(msg, "You're already authorized.")
			msg = append(msg, "Type /start to begin.")
			util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))
			return
		}
		newUser := users.User{
			ID:        m.Sender.ID,
			Username:  m.Sender.Username,
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
			Access:    users.UAMember,
		}
		b.Users.Create(newUser)

		// Notify User
		msg = append(msg, "You have been authorized.")
		msg = append(msg, "Type /start to begin.")
		util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))

		// Notify Admin
		adminMsg := fmt.Sprintf("%s has been granted acccess.", util.DisplayName(m.Sender))
		util.SendAdmin(b.TClient, b.Users.Admins(), adminMsg)
		return
	}

	// Notify User
	util.SendError(b.TClient, m.Sender, "Your password is invalid.")

	// Notify Admin
	adminMsg := "%s made an invalid auth request with password: %s"
	adminMsg = fmt.Sprintf(adminMsg, util.DisplayName(m.Sender), util.EscapeMarkdown(m.Payload))
	util.SendAdmin(b.TClient, b.Users.Admins(), adminMsg)
}
