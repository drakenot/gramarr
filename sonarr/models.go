package sonarr

import (
	"fmt"
	"time"
)

type TVShow struct {
	Title     string          `json:"title"`
	TitleSlug string          `json:"titleSlug"`
	Year      int             `json:"year"`
	PosterURL string          `json:"remotePoster"`
	TVDBID    int             `json:"tvdbId"`
	Images    []TVShowImage   `json:"images"`
	Seasons   []*TVShowSeason `json:"seasons"`
}

func (m TVShow) String() string {
	if m.Year != 0 {
		return fmt.Sprintf("%s (%d)", m.Title, m.Year)
	} else {
		return m.Title
	}
}

type TVShowImage struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
}

type TVShowSeason struct {
	SeasonNumber int  `json:"seasonNumber"`
	Monitored    bool `json:"monitored"`
}

type Folder struct {
	Path      string `json:"path"`
	FreeSpace int64  `json:"freeSpace"`
	ID        int    `json:"id"`
}

type Profile struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type AddTVShowRequest struct {
	Title             string           `json:"title"`
	TitleSlug         string           `json:"titleSlug"`
	Images            []TVShowImage    `json:"images"`
	QualityProfileID  int              `json:"qualityProfileId"`
	LanguageProfileID int              `json:"languageProfileId"`
	TVDBID            int              `json:"tvdbId"`
	RootFolderPath    string           `json:"rootFolderPath"`
	Monitored         bool             `json:"monitored"`
	AddOptions        AddTVShowOptions `json:"addOptions"`
	Year              int              `json:"year"`
	Seasons           []*TVShowSeason  `json:"seasons"`
}

type AddTVShowOptions struct {
	SearchForMissingEpisodes   bool `json:"searchForMissingEpisodes"`
	IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
}

type SystemStatus struct {
	Version           string    `json:"version"`
	BuildTime         time.Time `json:"buildTime"`
	IsDebug           bool      `json:"isDebug"`
	IsProduction      bool      `json:"isProduction"`
	IsAdmin           bool      `json:"isAdmin"`
	IsUserInteractive bool      `json:"isUserInteractive"`
	StartupPath       string    `json:"startupPath"`
	AppData           string    `json:"appData"`
	OsVersion         string    `json:"osVersion"`
	IsMono            bool      `json:"isMono"`
	IsLinux           bool      `json:"isLinux"`
	IsWindows         bool      `json:"isWindows"`
	Branch            string    `json:"branch"`
	Authentication    bool      `json:"authentication"`
	StartOfWeek       int       `json:"startOfWeek"`
	UrlBase           string    `json:"urlBase"`
}
