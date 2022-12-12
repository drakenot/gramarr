package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/pkg/chatbot"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) HandleCancel(m *tb.Message) {
	mb.tele.Send(m.Sender, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...")
}

func (mb *MediaBot) HandleConvoCancel(c chatbot.Conversation, m *tb.Message) {
	mb.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "util.Send /help for a list of commands.")
	mb.tele.Send(m.Sender, strings.Join(msg, "\n"))
}
