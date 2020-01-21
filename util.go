package main

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Send func
func Send(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

// SendError func
func SendError(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

// SendAdmin func
func SendAdmin(bot *tb.Bot, to []User, msg string) {
	SendMany(bot, to, fmt.Sprintf("*[Admin]* %s", msg))
}

// SendKeyboardList func
func SendKeyboardList(bot *tb.Bot, to tb.Recipient, msg string, list []string) {
	var buttons []tb.ReplyButton
	for _, item := range list {
		buttons = append(buttons, tb.ReplyButton{Text: item})
	}

	var replyKeys [][]tb.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []tb.ReplyButton{b})
	}

	bot.Send(to, msg, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
	})
}

// SendMany func
func SendMany(bot *tb.Bot, to []User, msg string) {
	for _, user := range to {
		bot.Send(user, msg, tb.ModeMarkdown)
	}
}

// DisplayName func
func DisplayName(u *tb.User) string {
	if u.FirstName != "" && u.LastName != "" {
		return EscapeMarkdown(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
	}

	return EscapeMarkdown(u.FirstName)
}

// EscapeMarkdown func
func EscapeMarkdown(s string) string {
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	return s
}
