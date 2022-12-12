package mediabot

import (
	"github.com/drakenot/gramarr/internal/repos/config"
	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/pkg/chatbot"
	"github.com/drakenot/gramarr/pkg/radarr"
	"github.com/drakenot/gramarr/pkg/sonarr"
	"github.com/drakenot/gramarr/pkg/telegram"
)

type MediaBot struct {
	*chatbot.ChatBot

	config *config.Config
	users  *users.UserDB
	tele   *telegram.Client
	radarr *radarr.Client
	sonarr *sonarr.Client
}

func NewMediaBot(config *config.Config, userDB *users.UserDB, tele *telegram.Client, radarr *radarr.Client, sonarr *sonarr.Client) *MediaBot {

	bot := &MediaBot{
		ChatBot: chatbot.NewChatBot(),
		config:  config,
		users:   userDB,
		tele:    tele,
		radarr:  radarr,
		sonarr:  sonarr,
	}

	// Send all telegram messages to our bot message router
	bot.tele.OnMessage(bot.Route)

	// Register commands with the bot
	bot.RegisterCommand("/auth", bot.RequirePrivate(bot.RequireAuth(users.UANone, bot.HandleAuth)))
	bot.RegisterCommand("/start", bot.RequirePrivate(bot.RequireAuth(users.UANone, bot.HandleStart)))
	bot.RegisterCommand("/help", bot.RequirePrivate(bot.RequireAuth(users.UANone, bot.HandleStart)))
	bot.RegisterCommand("/cancel", bot.RequirePrivate(bot.RequireAuth(users.UANone, bot.HandleCancel)))
	bot.RegisterCommand("/addmovie", bot.RequirePrivate(bot.RequireAuth(users.UAMember, bot.HandleAddMovie)))
	bot.RegisterCommand("/addtv", bot.RequirePrivate(bot.RequireAuth(users.UAMember, bot.HandleAddTVShow)))
	bot.RegisterCommand("/users", bot.RequirePrivate(bot.RequireAuth(users.UAAdmin, bot.HandleUsers)))

	// Catchall Command
	bot.RegisterFallback(bot.RequirePrivate(bot.RequireAuth(users.UANone, bot.HandleFallback)))

	// Conversation Commands
	bot.RegisterConvoCommand("/cancel", bot.HandleConvoCancel)

	return bot
}

func (mb *MediaBot) Start() {
	// Start Receiving telegram messages
	mb.tele.Start()
}
