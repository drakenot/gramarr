package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

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

func BoolToYesOrNo(condition bool) string {
	if condition {
		return "Yes"
	}
	return "No"
}

func FormatDate(dateStr string) string {
	if dateStr == "" {
		return "Unknown"
	}
	dateStr = strings.Split(dateStr, "T")[0]
	t, _ := time.Parse("2006-01-02", dateStr)
	return t.Format("02.01.2006")
}

func FormatDateTime(dateStr string) string {
	if dateStr == "" {
		return "Unknown"
	}
	dateStr = strings.Split(strings.Split(dateStr, ".")[0], "Z")[0]
	t, _ := time.Parse("2006-01-02T15:04:05", dateStr)
	return t.Format("02.01.2006 15:04:05")
}

func GetRootFolderFromPath(path string) string {
	return strings.Title(filepath.Base(filepath.Dir(path)))
}
