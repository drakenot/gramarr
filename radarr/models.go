package radarr

import "fmt"

type Movie struct {
	Title     string       `json:"title"`
	TitleSlug string       `json:"titleSlug"`
	Year      int          `json:"year"`
	PosterURL string       `json:"remotePoster"`
	TMDBID    int          `json:"tmdbId"`
	Images    []MovieImage `json:"images"`
}

func (m Movie) String() string {
	if m.Year != 0 {
		return fmt.Sprintf("%s (%d)", m.Title, m.Year)
	} else {
		return m.Title
	}
}

type MovieImage struct {
	CoverType string `json:"coverType"`
	URL       string `json:"url"`
}

type Folder struct {
	Path      string `json:"path"`
	FreeSpace int64  `json:"freeSpace"`
	ID        int    `json:"id"`
}

type AddMovieRequest struct {
	Title            string          `json:"title"`
	TitleSlug        string          `json:"titleSlug"`
	Images           []MovieImage    `json:"images"`
	QualityProfileID int             `json:"qualityProfileId"`
	TMDBID           int             `json:"tmdbId"`
	RootFolderPath   string          `json:"rootFolderPath"`
	Monitored        bool            `json:"monitored"`
	AddOptions       AddMovieOptions `json:"addOptions"`
}

type AddMovieOptions struct {
	SearchForMovie bool `json:"searchForMovie"`
}
