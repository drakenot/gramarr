package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

// Conversation interface
type Conversation interface {
	Run(m *tb.Message)
	CurrentStep() Handler
	Name() string
}

// ConversationManager struct
type ConversationManager struct {
	convos *cache.Cache
}

func NewConversationManager() *ConversationManager {
	convos := cache.New(30*time.Minute, 10*time.Minute)
	return &ConversationManager{convos: convos}
}

func (cm *ConversationManager) ProcessMessage(m *tb.Message) {
	key := cm.convoKey(m)
	if convo, ok := cm.convos.Get(key); ok {
		c := convo.(Conversation)
		c.CurrentStep()(m)
	}
}

func (cm *ConversationManager) HasConversation(m *tb.Message) bool {
	_, exists := cm.convos.Get(cm.convoKey(m))
	return exists
}

func (cm *ConversationManager) StartConversation(c Conversation, m *tb.Message) {
	c.Run(m)
	cm.convos.SetDefault(cm.convoKey(m), c)
}

func (cm *ConversationManager) StopConversation(c Conversation) {
	for key, item := range cm.convos.Items() {
		current := item.Object.(Conversation)
		if c == current {
			cm.convos.Delete(key)
		}
	}
}

func (cm *ConversationManager) Conversation(m *tb.Message) (Conversation, bool) {
	c, exists := cm.convos.Get(cm.convoKey(m))
	return c.(Conversation), exists
}

func (cm *ConversationManager) convoKey(m *tb.Message) string {
	return fmt.Sprintf("%d:%d", m.Chat.ID, m.Sender.ID)
}
