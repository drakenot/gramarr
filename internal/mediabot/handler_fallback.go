package mediabot

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) HandleFallback(m *tb.Message) {
	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	mb.tele.Send(m.Sender, strings.Join(msg, "\n"))
}
