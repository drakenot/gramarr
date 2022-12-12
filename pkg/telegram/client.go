package telegram

import (
	"fmt"
	Users "github.com/drakenot/gramarr/internal/repos/users"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

type Client struct {
	bot    *tb.Bot
	poller *tb.LongPoller
}

func NewClient(token string, pollTimeout time.Duration) (*Client, error) {
	poller := &tb.LongPoller{Timeout: pollTimeout}
	bot, err := tb.NewBot(tb.Settings{Token: token, Poller: poller})
	if err != nil {
		return nil, err
	}
	return &Client{
		poller: poller,
		bot:    bot,
	}, nil
}

func (c *Client) Start() {
	c.bot.Start()
}

func (c *Client) OnMessage(handler func(m *tb.Message)) {
	c.bot.Handle(tb.OnText, handler)
}

// Send will send a message to a user
func (c *Client) Send(to tb.Recipient, msg interface{}) {
	switch msg.(type) {
	case string:
		c.bot.Send(to, msg, tb.ModeMarkdown)
		break
	default:
		c.bot.Send(to, msg)
	}
}

func (c *Client) Me() *tb.User {
	return c.bot.Me
}

// SendError will send an error message to a user
func (c *Client) SendError(to tb.Recipient, msg string) {
	c.bot.Send(to, msg, tb.ModeMarkdown)
}

// SendAdmin will send an error to all admin users
func (c *Client) SendAdmin(to []Users.User, msg string) {
	c.SendMany(to, fmt.Sprintf("*[Admin]* %s", msg))
}

// SendChoices will send a list of keyboard choices to a user
func (c *Client) SendChoices(to tb.Recipient, msg string, list []string) {
	var buttons []tb.ReplyButton
	for _, item := range list {
		buttons = append(buttons, tb.ReplyButton{Text: item})
	}

	var replyKeys [][]tb.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []tb.ReplyButton{b})
	}

	c.bot.Send(to, msg, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
	})
}

// SendMany will send a message to multiple users in markdown format
func (c *Client) SendMany(to []Users.User, msg string) {
	for _, user := range to {
		c.bot.Send(user, msg, tb.ModeMarkdown)
	}
}

// DisplayName will generate a "display name" for a user
func (c *Client) DisplayName(u *tb.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return c.escapeMarkdown(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	}

	return c.escapeMarkdown(u.FirstName)
}

// GetUserName will generate a "user name" for a given user
func (c *Client) GetUserName(m *tb.Message) string {
	var username string
	if len(m.Sender.Username) > -1 {
		username = m.Sender.Username
	} else {
		username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
	}
	return strings.TrimSpace(strings.ToLower(username))
}

// EscapeMarkdown will escape any markdown characters in a string
func (c *Client) escapeMarkdown(s string) string {
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return s
}
