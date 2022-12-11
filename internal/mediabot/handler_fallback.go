package mediabot

import (
	"strings"

	"github.com/drakenot/gramarr/internal/util"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *MediaBot) HandleFallback(m *tb.Message) {
	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))
}
