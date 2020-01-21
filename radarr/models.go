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

type Profile struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Folder struct {
	Path      string `json:"path"`
	FreeSpace int64  `json:"freeSpace"`
	ID        int    `json:"id"`
}

type AddMovieRequest struct {
	Title             string          `json:"title"`
	TitleSlug         string          `json:"titleSlug"`
	Images            []MovieImage    `json:"images"`
	QualityProfileID  int             `json:"qualityProfileId"`
	LanguageProfileID int             `json:"languageProfileId"`
	TMDBID            int             `json:"tmdbId"`
	RootFolderPath    string          `json:"rootFolderPath"`
	Monitored         bool            `json:"monitored"`
	AddOptions        AddMovieOptions `json:"addOptions"`
	Year              int             `json:"year"`
}

type AddMovieOptions struct {
	SearchForMovie bool `json:"searchForMovie"`
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
