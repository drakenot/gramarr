package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/drakenot/gramarr/internal/mediabot"
	"github.com/drakenot/gramarr/internal/repos/config"
	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/pkg/chatbot"
	"github.com/drakenot/gramarr/pkg/radarr"
	"github.com/drakenot/gramarr/pkg/sonarr"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Flags
var (
	configDir = flag.String("configDir", ".", "config dir for settings and logs")
)

func main() {
	flag.Parse()

	conf, err := config.LoadConfig(*configDir)
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	err = config.ValidateConfig(conf)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	userPath := filepath.Join(*configDir, "users.json")
	users, err := users.NewUserDB(userPath)
	if err != nil {
		log.Fatalf("failed to load the user db %v", err)
	}

	var rc *radarr.Client
	if conf.Radarr != nil {
		rc, err = radarr.NewClient(*conf.Radarr)
		if err != nil {
			log.Fatalf("failed to create radarr client: %v", err)
		}
	}

	var sn *sonarr.Client
	if conf.Sonarr != nil {
		sn, err = sonarr.NewClient(*conf.Sonarr)
		if err != nil {
			log.Fatalf("failed to create sonarr client: %v", err)
		}
	}

	cm := chatbot.NewConversationManager()
	router := chatbot.NewRouter(cm)

	poller := tb.LongPoller{Timeout: 15 * time.Second}
	bot, err := tb.NewBot(tb.Settings{
		Token:  conf.Telegram.BotToken,
		Poller: &poller,
	})
	if err != nil {
		log.Fatalf("failed to create telegram bot client: %v", err)
	}

	env := &mediabot.MediaBot{
		Config:  conf,
		TClient: bot,
		Users:   users,
		CM:      cm,
		Radarr:  rc,
		Sonarr:  sn,
	}

	setupHandlers(router, env)
	fmt.Fprintf(os.Stdout, "Gramarr is up and running. Go call your bot!\n")
	bot.Start()
}

func setupHandlers(r *chatbot.Router, b *mediabot.MediaBot) {
	// Send all telegram messages to our custom router
	b.TClient.Handle(tb.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", b.RequirePrivate(b.RequireAuth(users.UANone, b.HandleAuth)))
	r.HandleFunc("/start", b.RequirePrivate(b.RequireAuth(users.UANone, b.HandleStart)))
	r.HandleFunc("/help", b.RequirePrivate(b.RequireAuth(users.UANone, b.HandleStart)))
	r.HandleFunc("/cancel", b.RequirePrivate(b.RequireAuth(users.UANone, b.HandleCancel)))
	r.HandleFunc("/addmovie", b.RequirePrivate(b.RequireAuth(users.UAMember, b.HandleAddMovie)))
	r.HandleFunc("/addtv", b.RequirePrivate(b.RequireAuth(users.UAMember, b.HandleAddTVShow)))
	r.HandleFunc("/users", b.RequirePrivate(b.RequireAuth(users.UAAdmin, b.HandleUsers)))

	// Catchall Command
	r.HandleFallback(b.RequirePrivate(b.RequireAuth(users.UANone, b.HandleFallback)))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", b.HandleConvoCancel)
}
