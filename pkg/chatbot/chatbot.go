package chatbot

import (
	"fmt"
	"regexp"
	"time"

	"github.com/patrickmn/go-cache"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	cmdRx = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
)

type Handler func(*tb.Message)
type ConvoHandler func(Conversation, *tb.Message)

type Conversation interface {
	Run(m *tb.Message)
	CurrentStep() Handler
	Name() string
}

type ChatBot struct {
	convos      *cache.Cache
	routes      map[string]Handler
	convoRoutes map[string]ConvoHandler
	fallback    Handler
}

func NewChatBot() *ChatBot {
	convos := cache.New(30*time.Minute, 10*time.Minute)
	return &ChatBot{convos: convos,
		routes:      map[string]Handler{},
		convoRoutes: map[string]ConvoHandler{},
	}
}

func (cb *ChatBot) ProcessMessage(m *tb.Message) {
	key := cb.convoKey(m)
	if convo, ok := cb.convos.Get(key); ok {
		c := convo.(Conversation)
		c.CurrentStep()(m)
	}
}

func (cb *ChatBot) HasConversation(m *tb.Message) bool {
	_, exists := cb.convos.Get(cb.convoKey(m))
	return exists
}

func (cb *ChatBot) StartConversation(c Conversation, m *tb.Message) {
	c.Run(m)
	cb.convos.SetDefault(cb.convoKey(m), c)
}

func (cb *ChatBot) StopConversation(c Conversation) {
	for key, item := range cb.convos.Items() {
		current := item.Object.(Conversation)
		if c == current {
			cb.convos.Delete(key)
		}
	}
}

func (cb *ChatBot) Conversation(m *tb.Message) (Conversation, bool) {
	c, exists := cb.convos.Get(cb.convoKey(m))
	return c.(Conversation), exists
}
func (cb *ChatBot) RegisterCommand(cmd string, h Handler) {
	cb.routes[cmd] = h
}

func (cb *ChatBot) RegisterFallback(h Handler) {
	cb.fallback = h
}

func (cb *ChatBot) RegisterConvoCommand(cmd string, h ConvoHandler) {
	cb.convoRoutes[cmd] = h
}

func (cb *ChatBot) Route(m *tb.Message) {
	if !cb.routeConvo(m) && !cb.routeCommand(m) {
		cb.routeFallback(m)
	}
}

func (cb *ChatBot) routeConvo(m *tb.Message) bool {
	if !cb.HasConversation(m) {
		return false
	}

	// Global Conversation Cmd?
	if cmd, match := cb.parseCommand(m); match {
		if route, exists := cb.convoRoutes[cmd]; exists {
			convo, _ := cb.Conversation(m)
			route(convo, m)
			return true
		}
	}

	cb.ProcessMessage(m)
	return true
}

func (cb *ChatBot) routeCommand(m *tb.Message) bool {
	if cmd, match := cb.parseCommand(m); match {
		if route, exists := cb.routes[cmd]; exists {
			route(m)
			return true
		}
	}
	return false
}

func (cb *ChatBot) routeFallback(m *tb.Message) {
	if cb.fallback != nil {
		cb.fallback(m)
	}
}

func (cb *ChatBot) parseCommand(m *tb.Message) (string, bool) {
	match := cmdRx.FindAllStringSubmatch(m.Text, -2)

	if match != nil {
		return match[0][1], true
	} else {
		return "", false
	}
}

func (cb *ChatBot) convoKey(m *tb.Message) string {
	return fmt.Sprintf("%d:%d", m.Chat.ID, m.Sender.ID)
}
