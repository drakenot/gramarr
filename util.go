package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func Send(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

func SendError(bot *tb.Bot, to tb.Recipient, msg string) {
	bot.Send(to, msg, tb.ModeMarkdown)
}

func SendAdmin(bot *tb.Bot, to []User, msg string) {
	SendMany(bot, to, fmt.Sprintf("*[Admin]* %s", msg))
}

func SendKeyboardList(bot *tb.Bot, to tb.Recipient, msg string, list []string) {
	var buttons []tb.ReplyButton
	for _, item := range list {
		buttons = append(buttons, tb.ReplyButton{Text:item})
	}

	var replyKeys [][]tb.ReplyButton
	for _, b := range buttons {
		replyKeys = append(replyKeys, []tb.ReplyButton{b})
	}

	bot.Send(to, msg, &tb.ReplyMarkup{
		ReplyKeyboard:  replyKeys,
		OneTimeKeyboard:true,
	})
}

func SendMany(bot *tb.Bot, to []User, msg string) {
	for _, user := range to {
		bot.Send(user, msg, tb.ModeMarkdown)
	}
}

func DisplayName(u *tb.User) string {
	if u.Username != "" {
		return u.Username
	} else {
		if u.LastName != "" {
			return u.FirstName + " " + u.LastName
		}
		return u.FirstName
	}
}
