package util

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// DisplayName will generate a "display name" for a user
func DisplayName(u *tb.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return EscapeMarkdown(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	}

	return EscapeMarkdown(u.FirstName)
}

// EscapeMarkdown will escape any markdown characters in a string
func EscapeMarkdown(s string) string {
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return s
}

// GetUserName will generate a "user name" for a given user
func GetUserName(m *tb.Message) string {
	var username string
	if len(m.Sender.Username) > 0 {
		username = m.Sender.Username
	} else {
		username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
	}
	return strings.TrimSpace(strings.ToLower(username))
}
