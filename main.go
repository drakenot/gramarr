package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gramarr/radarr"
	"github.com/gramarr/sonarr"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	configDir = flag.String("configDir", "./config", "config dir for settings and logs")
)

type Env struct {
	Config *Config
	Users  *UserDB
	Bot    *tb.Bot
	CM     *ConversationManager
	Radarr *radarr.Client
	Sonarr *sonarr.Client
}

func main() {
	flag.Parse()

	conf, err := LoadConfig(*configDir)
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}

	err = ValidateConfig(conf)
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	userPath := filepath.Join(*configDir, "users.json")
	users, err := NewUserDB(userPath)
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
		Sonarr: sn,
	}

	setupHandlers(router, env)
	fmt.Fprintf(os.Stdout, "Gramarr is up and running. Go call your bot!\n")
	bot.Start()
}

func setupHandlers(r *Router, e *Env) {
	// Send all telegram messages to our custom router
	e.Bot.Handle(tb.OnText, r.Route)

	// Commands
	r.HandleFunc("/auth", e.RequirePrivate(e.RequireAuth(UANone, e.HandleAuth)))
	r.HandleFunc("/start", e.RequirePrivate(e.RequireAuth(UANone, e.HandleHelp)))
	r.HandleFunc("/help", e.RequirePrivate(e.RequireAuth(UANone, e.HandleHelp)))
	r.HandleFunc("/cancel", e.RequirePrivate(e.RequireAuth(UANone, e.HandleCancel)))
	r.HandleFunc("/addmovie", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleAddMovie)))
	r.HandleFunc("/listmovies", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleListMovies)))
	r.HandleFunc("/moviedetails", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleMovieDetails)))
	r.HandleFunc("/deletemovie", e.RequirePrivate(e.RequireAuth(UAAdmin, e.HandleDeleteMovie)))
	r.HandleFunc("/addtv", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleAddTVShow)))
	r.HandleFunc("/listtv", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleListTVShows)))
	r.HandleFunc("/tvdetails", e.RequirePrivate(e.RequireAuth(UAMember, e.HandleTVShowDetails)))
	r.HandleFunc("/deletetv", e.RequirePrivate(e.RequireAuth(UAAdmin, e.HandleDeleteTVShow)))
	r.HandleFunc("/users", e.RequirePrivate(e.RequireAuth(UAAdmin, e.HandleUsers)))
	r.HandleFunc("/status", e.RequirePrivate(e.RequireAuth(UAAdmin, e.HandleStatus)))

	// Catchall Command
	r.HandleFallback(e.RequirePrivate(e.RequireAuth(UANone, e.HandleFallback)))

	// Conversation Commands
	r.HandleConvoFunc("/cancel", e.HandleConvoCancel)
}
