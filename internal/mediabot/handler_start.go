package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *MediaBot) HandleStart(m *tb.Message) {

	user, exists := b.Users.User(m.Sender.ID)

	var msg []string
	msg = append(msg, fmt.Sprintf("Hello, I'm %s! Use these commands to control me:", b.TClient.Me.FirstName))

	if !exists {
		msg = append(msg, "")
		msg = append(msg, "/auth [password] - authenticate with the bot")
	}

	if exists && user.IsAdmin() {
		msg = append(msg, "")
		msg = append(msg, "*Admin*")
		msg = append(msg, "/users - list all bot users")
	}

	if exists && (user.IsMember() || user.IsAdmin()) {
		msg = append(msg, "")
		msg = append(msg, "*Media*")
		msg = append(msg, "/addmovie - add a movie")
		msg = append(msg, "/addtv - add a tv show")
		msg = append(msg, "")
		msg = append(msg, "/cancel - cancel the current operation")
	}

	util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))
}
