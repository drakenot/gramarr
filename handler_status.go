package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"reflect"
	"strings"
)

func (e *Env) HandleStatus(m *tb.Message) {

	sonarrStatus, err := e.Sonarr.GetSystemStatus()
	if err != nil {
		SendError(e.Bot, m.Sender, "Failed to get Sonarr System Status.")
	} else {
		var msg []string
		msg = append(msg, "*Sonarr Status:*")
		s := reflect.ValueOf(&sonarrStatus).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			msg = append(msg, fmt.Sprintf("%s = %v", typeOfT.Field(i).Name, f.Interface()))
		}
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
	}

	radarrStatus, err := e.Radarr.GetSystemStatus()
	if err != nil {
		SendError(e.Bot, m.Sender, "Failed to get Radarr System Status.")
	} else {
		var msg []string
		msg = append(msg, "*Radarr Status:*")
		s := reflect.ValueOf(&radarrStatus).Elem()
		typeOfT := s.Type()
		for i := 0; i < s.NumField(); i++ {
			f := s.Field(i)
			msg = append(msg, fmt.Sprintf("%s = %v", typeOfT.Field(i).Name, f.Interface()))
		}
		Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
	}

}
