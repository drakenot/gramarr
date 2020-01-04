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

// NewConversationManager func
func NewConversationManager() *ConversationManager {
	convos := cache.New(30*time.Minute, 10*time.Minute)
	return &ConversationManager{convos: convos}
}

// ProcessMessage func
func (cm *ConversationManager) ProcessMessage(m *tb.Message) {
	key := cm.convoKey(m)
	if convo, ok := cm.convos.Get(key); ok {
		c := convo.(Conversation)
		c.CurrentStep()(m)
	}
}

// HasConversation func
func (cm *ConversationManager) HasConversation(m *tb.Message) bool {
	_, exists := cm.convos.Get(cm.convoKey(m))
	return exists
}

// StartConversation func
func (cm *ConversationManager) StartConversation(c Conversation, m *tb.Message) {
	c.Run(m)
	cm.convos.SetDefault(cm.convoKey(m), c)
}

// StopConversation func
func (cm *ConversationManager) StopConversation(c Conversation) {
	for key, item := range cm.convos.Items() {
		current := item.Object.(Conversation)
		if c == current {
			cm.convos.Delete(key)
		}
	}
}

// Conversation func
func (cm *ConversationManager) Conversation(m *tb.Message) (Conversation, bool) {
	c, exists := cm.convos.Get(cm.convoKey(m))
	return c.(Conversation), exists
}

// convoKey func
func (cm *ConversationManager) convoKey(m *tb.Message) string {
	return fmt.Sprintf("%d:%d", m.Chat.ID, m.Sender.ID)
}
