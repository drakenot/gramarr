package main

import (
	"regexp"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	cmdRx = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
)

// Handler func
type Handler func(*tb.Message)

// ConvoHandler func
type ConvoHandler func(Conversation, *tb.Message)

// NewRouter func
func NewRouter(cm *ConversationManager) *Router {
	return &Router{cm: cm, routes: map[string]Handler{}, convoRoutes: map[string]ConvoHandler{}}
}

// Router struct
type Router struct {
	cm          *ConversationManager
	routes      map[string]Handler
	convoRoutes map[string]ConvoHandler
	fallback    Handler
}

// HandleFunc func
func (r *Router) HandleFunc(cmd string, h Handler) {
	r.routes[cmd] = h
}

// HandleFallback func
func (r *Router) HandleFallback(h Handler) {
	r.fallback = h
}

// HandleConvoFunc func
func (r *Router) HandleConvoFunc(cmd string, h ConvoHandler) {
	r.convoRoutes[cmd] = h
}

// Route func
func (r *Router) Route(m *tb.Message) {
	if !r.routeConvo(m) && !r.routeCommand(m) {
		r.routeFallback(m)
	}
}

// routeConvo func
func (r *Router) routeConvo(m *tb.Message) bool {
	if !r.cm.HasConversation(m) {
		return false
	}

	// Global Conversation Cmd?
	if cmd, match := r.parseCommand(m); match {
		if route, exists := r.convoRoutes[cmd]; exists {
			convo, _ := r.cm.Conversation(m)
			route(convo, m)
			return true
		}
	}

	r.cm.ProcessMessage(m)
	return true
}

// routeCommand func
func (r *Router) routeCommand(m *tb.Message) bool {
	if cmd, match := r.parseCommand(m); match {
		if route, exists := r.routes[cmd]; exists {
			route(m)
			return true
		}
	}
	return false
}

// routeFallback func
func (r Router) routeFallback(m *tb.Message) {
	if r.fallback != nil {
		r.fallback(m)
	}
}

// parseCommand func
func (r *Router) parseCommand(m *tb.Message) (string, bool) {
	match := cmdRx.FindAllStringSubmatch(m.Text, -1)

	if match != nil {
		return match[0][1], true
	} else {
		return "", false
	}
}
