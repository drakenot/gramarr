package sonarr

import "fmt"

type TVShow struct {
	Title     string         `json:"title"`
	TitleSlug string         `json:"titleSlug"`
	Year      int            `json:"year"`
	PosterURL string         `json:"remotePoster"`
	TVDBID    int            `json:"tvdbId"`
	Images    []TVShowImage  `json:"images"`
	Seasons   []TVShowSeason `json:"seasons"`
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
	Seasons           []TVShowSeason   `json:"seasons"`
}

type AddTVShowOptions struct {
	SearchForMissingEpisodes   bool `json:"searchForMissingEpisodes"`
	IgnoreEpisodesWithFiles    bool `json:"ignoreEpisodesWithFiles"`
	IgnoreEpisodesWithoutFiles bool `json:"ignoreEpisodesWithoutFiles"`
}
