package sonarr

import "fmt"

type TVShow struct {
	Title             string                 `json:"title"`
	TitleSlug         string                 `json:"titleSlug"`
	Year              int                    `json:"year"`
	PosterURL         string                 `json:"remotePoster"`
	TVDBID            int                    `json:"tvdbId"`
	Images            []TVShowImage          `json:"images"`
	Seasons           []*TVShowSeason        `json:"seasons"`
	AlternateTitles   []TVShowAlternateTitle `json:"alternateTitles"`
	SortTitle         string                 `json:"sortTitle"`
	SeasonCount       int                    `json:"seasonCount"`
	TotalEpisodeCount int                    `json:"TotalEpisodeCount"`
	EpisodeCount      int                    `json:"episodeCount"`
	EpisodeFileCount  int                    `json:"episodeFileCount"`
	SizeOnDisk        uint64                 `json:"sizeOnDisk"`
	Status            string                 `json:"status"`
	Overview          string                 `json:"overview"`
	PreviousAiring    string                 `json:"previousAiring"`
	Network           string                 `json:"network"`
	AirTime           string                 `json:"airTime"`
	Path              string                 `json:"path"`
	ProfileID         int                    `json:"profileId"`
	SeasonFolder      bool                   `json:"seasonFolder"`
	Monitored         bool                   `json:"monitored"`
	UseSceneNumbering bool                   `json:"useSceneNumbering"`
	Runtime           int                    `json:"runtime"`
	TVRageID          int                    `json:"tvRageId"`
	TVMazeID          int                    `json:"tvMazeId"`
	FirstAired        string                 `json:"firstAired"`
	LastInfoSync      string                 `json:"lastInfoSync"`
	SeriesType        string                 `json:"seriesType"`
	CleanTitle        string                 `json:"cleanTitle"`
	IMDBID            string                 `json:"imdbId"`
	Certification     string                 `json:"certification"`
	Genres            []string               `json:"genres"`
	Tags              []string               `json:"tags"`
	Added             string                 `json:"added"`
	Ratings           TVShowRating           `json:"ratings"`
	QualityProfileID  int                    `json:"qualityProfileId"`
	ID                int                    `json:"id"`
}

type TVShowAlternateTitle struct {
	Title        string `json:"title"`
	SeasonNumber int    `json:"seasonNumber"`
}

type TVShowRating struct {
	Votes int     `json:"votes"`
	Value float32 `json:"value"`
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
	SeasonFolder      bool             `json:"seasonFolder"`
}

type AddTVShowOptions struct {
	SearchForMissingEpisodes   bool `json:"searchForMissingEpisodes"`
	IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
}
