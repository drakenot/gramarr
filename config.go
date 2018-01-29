package main

import (
	"encoding/json"
	"github.com/drakenot/gramarr/radarr"
	"github.com/drakenot/gramarr/sonarr"
	"io/ioutil"
	"log"
)

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	Bot      BotConfig      `json:"bot"`
	Radarr   *radarr.Config `json:"radarr"`
	Sonarr   *sonarr.Config `json:"sonarr"`
}

type TelegramConfig struct {
	BotToken string `json:"botToken"`
}

type BotConfig struct {
	UserDBPath    string `json:"userDbPath"`
	Password      string `json:"password"`
	AdminPassword string `json:"adminPassword"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	var c Config
	json.Unmarshal(file, &c)
	return &c, nil
}
