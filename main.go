package main

import (
	"flag"
	"log"
	"time"

	"github.com/drakenot/gramarr/radarr"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Flags
var (
	configPath = flag.String("config", "config.json", "config.json path")
)

type Env struct {
	Config *Config
	Users  *UserDB
	Bot    *tb.Bot
	CM     *ConversationManager
	Radarr *radarr.Client
}

func main() {
	flag.Parse()

	conf, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config file at %s: %v", *configPath, err)
	}

	err = validateConfig(conf)
	if err != nil {
		log.Fatal("config error: %v", err)
	}

	users, err := NewUserDB(conf.Bot.UserDBPath)
	if err != nil {
		log.Fatalf("failed to load the user db at %s: %v", conf.Bot.UserDBPath, err)
	}

	var rc *radarr.Client
	if conf.Radarr != nil {
		rc, err = radarr.NewClient(*conf.Radarr)
		if err != nil {
			log.Fatalf("failed to create radarr client: %v", err)
		}
	}

	cm := NewConversationManager()
	router := NewRouter(cm)

	poller := tb.LongPoller{Timeout: 15 * time.Second}
	bot, err := tb.NewBot(tb.Settings{
		Token:  conf.Telegram.BotToken,
		Poller: &poller,
	})
	if err != nil {
		log.Fatalf("failed to create telegram bot client: %v", err)
	}

	env := &Env{
		Config: conf,
		Bot:    bot,
		Users:  users,
		CM:     cm,
		Radarr: rc,
	}

	setupHandlers(router, env)
	bot.Start()
}

func setupHandlers(r *Router, e *Env) {
	// Send all telegram messages to our custom router
	e.Bot.Handle(tb.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", e.RequirePrivate(e.RequireAuth(UANone, e.HandleAuth)))
	r.HandleFunc("/start", e.RequirePrivate(e.RequireAuth(UANone, e.HandleStart)))
	r.HandleFunc("/help", e.RequirePrivate(e.RequireAuth(UANone, e.HandleStart)))
	r.HandleFunc("/cancel", e.RequirePrivate(e.RequireAuth(UANone, e.HandleCancel)))
	r.HandleFunc("/addmovie", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleAddMovie)))
	r.HandleFunc("/users", e.RequirePrivate(e.RequireAuth(UAAdmin, e.HandleUsers)))

	// Catchall Command
	r.HandleFallback(e.RequirePrivate(e.RequireAuth(UANone, e.HandleFallback)))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", e.HandleConvoCancel)
}

func validateConfig(c *Config) error {
	return nil
}
