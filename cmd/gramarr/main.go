// Package is Media Bot application that interfaces with sonarr/radarr to manage media
package main

import (
	"flag"
	"fmt"
	"github.com/drakenot/gramarr/internal/mediabot"
	"github.com/drakenot/gramarr/internal/repos/config"
	"github.com/drakenot/gramarr/internal/repos/users"
	"github.com/drakenot/gramarr/pkg/radarr"
	"github.com/drakenot/gramarr/pkg/sonarr"
	"github.com/drakenot/gramarr/pkg/telegram"
	"log"
	"os"
	"path/filepath"
	"time"
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

	var sc *sonarr.Client
	if conf.Sonarr != nil {
		sc, err = sonarr.NewClient(*conf.Sonarr)
		if err != nil {
			log.Fatalf("failed to create sonarr client: %v", err)
		}
	}

	tc, err := telegram.NewClient(conf.Telegram.BotToken, time.Minute*14)
	if err != nil {
		log.Fatalf("failed to create telegram client: #{err}")
	}

	mb := mediabot.NewMediaBot(conf, users, tc, rc, sc)

	fmt.Fprintf(os.Stdout, "Gramarr is up and running. Go call your bot!\n")
	mb.Start()
}
