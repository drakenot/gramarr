package mediabot

import (
	"fmt"
	"strings"

	"github.com/drakenot/gramarr/internal/util"
	"github.com/drakenot/gramarr/pkg/chatbot"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *MediaBot) HandleCancel(m *tb.Message) {
	util.Send(b.TClient, m.Sender, "There is no active command to cancel. I wasn't doing anything anyway. Zzzzz...")
}

func (b *MediaBot) HandleConvoCancel(c chatbot.Conversation, m *tb.Message) {
	b.CM.StopConversation(c)

	var msg []string
	msg = append(msg, fmt.Sprintf("The '*%s*' command was cancelled. Anything else I can do for you?", c.Name()))
	msg = append(msg, "")
	msg = append(msg, "util.Send /help for a list of commands.")
	util.Send(b.TClient, m.Sender, strings.Join(msg, "\n"))
}
