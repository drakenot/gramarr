package mediabot

import (
	"github.com/drakenot/gramarr/internal/repos/config"
	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/pkg/chatbot"
	"github.com/drakenot/gramarr/pkg/radarr"
	"github.com/drakenot/gramarr/pkg/sonarr"

	tb "gopkg.in/tucnak/telebot.v2"
)

type MediaBot struct {
	Config  *config.Config
	Users   *users.UserDB
	TClient *tb.Bot
	CM      *chatbot.ConversationManager
	Radarr  *radarr.Client
	Sonarr  *sonarr.Client
}
