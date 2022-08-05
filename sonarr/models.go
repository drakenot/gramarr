package sonarr

import (
	"fmt"
	"time"
)

type TVShow struct {
	Added             time.Time         `json:"added"`
	AirTime           string            `json:"airTime,omitempty"`
	AlternateTitles   []AlternateTitles `json:"alternateTitles"`
	Certification     string            `json:"certification,omitempty"`
	CleanTitle        string            `json:"cleanTitle"`
	Ended             bool              `json:"ended"`
	FirstAired        time.Time         `json:"firstAired,omitempty"`
	Genres            []string          `json:"genres"`
	ID                int               `json:"id,omitempty"`
	Images            []TVShowImage     `json:"images"`
	ImdbID            string            `json:"imdbId,omitempty"`
	LanguageProfileID int               `json:"languageProfileId"`
	Monitored         bool              `json:"monitored,omitempty"`
	Network           string            `json:"network"`
	NextAiring        time.Time         `json:"nextAiring,omitempty"`
	Overview          string            `json:"overview,omitempty"`
	Path              string            `json:"path,omitempty"`
	PreviousAiring    time.Time         `json:"previousAiring,omitempty"`
	QualityProfileID  int               `json:"qualityProfileId"`
	Ratings           TVShowRatings     `json:"ratings,omitempty"`
	RootFolderPath    string            `json:"rootFolderPath"`
	Runtime           int               `json:"runtime"`
	RemotePoster      string            `json:"remotePoster"`
	Seasons           []*TVShowSeason   `json:"seasons"`
	SeasonFolder      bool              `json:"seasonFolder"`
	SeriesType        string            `json:"seriesType"`
	SortTitle         string            `json:"sortTitle"`
	Statistics        TVShowStatistics  `json:"statistics,omitempty"`
	Status            string            `json:"status"`
	Tags              []int             `json:"tags,omitempty"`
	Title             string            `json:"title"`
	TitleSlug         string            `json:"titleSlug"`
	TvMazeID          int               `json:"tvMazeId"`
	TvRageID          int               `json:"tvRageId"`
	TvdbID            int               `json:"tvdbId"`
	UseSceneNumbering bool              `json:"useSceneNumbering"`
	Year              int               `json:"year"`
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
	RemoteURL string `json:"remoteUrl,omitempty"`
}

type AlternateTitles struct {
	Title        string `json:"title,omitempty"`
	SeasonNumber int    `json:"seasonNumber,omitempty"`
}

type TVShowSeason struct {
	SeasonNumber int              `json:"seasonNumber"`
	Monitored    bool             `json:"monitored"`
	Statistics   SeasonStatistics `json:"statistics,omitempty"`
}

type TVShowRatings struct {
	Value float64 `json:"value,omitempty"`
	Votes int     `json:"votes,omitempty"`
}

type TVShowStatistics struct {
	SeasonCount       int     `json:"seasonCount,omitempty"`
	EpisodeFileCount  int     `json:"episodeFileCount,omitempty"`
	EpisodeCount      int     `json:"episodeCount,omitempty"`
	TotalEpisodeCount int     `json:"totalEpisodeCount,omitempty"`
	SizeOnDisk        int64   `json:"sizeOnDisk,omitempty"`
	PercentOfEpisodes float64 `json:"percentOfEpisodes,omitempty"`
}

type SeasonStatistics struct {
	EpisodeCount      int       `json:"episodeCount"`
	EpisodeFileCount  int       `json:"episodeFileCount"`
	PercentOfEpisodes float64   `json:"percentOfEpisodes"`
	NextAiring        time.Time `json:"nextAiring"`
	PreviousAiring    time.Time `json:"previousAiring"`
	SizeOnDisk        int64     `json:"sizeOnDisk"`
	TotalEpisodeCount int       `json:"totalEpisodeCount"`
}

type TVShowTag struct {
	Id    int    `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
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
	SeasonFolder      bool             `json:"seasonFolder"`
	Tags              []int            `json:"tags,omitempty"`
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
	Version           string `json:"version"`
	BuildTime         string `json:"buildTime"`
	IsDebug           bool   `json:"isDebug"`
	IsProduction      bool   `json:"isProduction"`
	IsAdmin           bool   `json:"isAdmin"`
	IsUserInteractive bool   `json:"isUserInteractive"`
	StartupPath       string `json:"startupPath"`
	AppData           string `json:"appData"`
	OsName            string `json:"osName"`
	OsVersion         string `json:"osVersion"`
	IsMonoRuntime     bool   `json:"isMonoRuntime"`
	IsMono            bool   `json:"isMono"`
	IsLinux           bool   `json:"isLinux"`
	IsOsx             bool   `json:"isOsx"`
	IsWindows         bool   `json:"isWindows"`
	Branch            string `json:"branch"`
	Authentication    string `json:"authentication"`
	SqliteVersion     string `json:"sqliteVersion"`
	UrlBase           string `json:"urlBase"`
	RuntimeVersion    string `json:"runtimeVersion"`
	RuntimeName       string `json:"runtimeName"`
}
