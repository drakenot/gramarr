package radarr

import (
	"fmt"
	"time"
)

type Movie struct {
	Title                 string             `json:"title,omitempty,omitempty"`
	AlternativeTitles     []AlternativeTitle `json:"alternativeTitles,omitempty"`
	SecondaryYearSourceID int                `json:"secondaryYearSourceId,omitempty"`
	SortTitle             string             `json:"sortTitle,omitempty"`
	SizeOnDisk            int64              `json:"sizeOnDisk,omitempty"`
	Status                string             `json:"status,omitempty"`
	Overview              string             `json:"overview,omitempty"`
	InCinemas             time.Time          `json:"inCinemas,omitempty"`
	PhysicalRelease       time.Time          `json:"physicalRelease,omitempty"`
	Images                []Image            `json:"images,omitempty"`
	Website               string             `json:"website,omitempty"`
	Downloaded            bool               `json:"downloaded,omitempty"`
	Year                  int                `json:"year,omitempty"`
	HasFile               bool               `json:"hasFile,omitempty"`
	YouTubeTrailerID      string             `json:"youTubeTrailerId,omitempty"`
	Studio                string             `json:"studio,omitempty"`
	Path                  string             `json:"path,omitempty"`
	ProfileID             int                `json:"profileId,omitempty"`
	Monitored             bool               `json:"monitored,omitempty"`
	MinimumAvailability   string             `json:"minimumAvailability,omitempty"`
	IsAvailable           bool               `json:"isAvailable,omitempty"`
	FolderName            string             `json:"folderName,omitempty"`
	Runtime               int                `json:"runtime,omitempty"`
	LastInfoSync          time.Time          `json:"lastInfoSync,omitempty"`
	CleanTitle            string             `json:"cleanTitle,omitempty"`
	ImdbID                string             `json:"imdbId,omitempty"`
	TmdbID                int                `json:"tmdbId,omitempty"`
	TitleSlug             string             `json:"titleSlug,omitempty"`
	Genres                []string           `json:"genres,omitempty"`
	Tags                  []int              `json:"tags,omitempty"`
	Added                 time.Time          `json:"added,omitempty"`
	Ratings               Ratings            `json:"ratings,omitempty"`
	QualityProfileID      int                `json:"qualityProfileId,omitempty"`
	ID                    int                `json:"id,omitempty"`
	MovieFile             MovieFile          `json:"movieFile,omitempty"`
	RemotePoster          string             `json:"remotePoster,omitempty"`
}

func (m Movie) String() string {
	if m.Year != 0 {
		return fmt.Sprintf("%s (%d)", m.Title, m.Year)
	} else {
		return m.Title
	}
}

type MovieFile struct {
	MovieID      int       `json:"movieId,omitempty"`
	RelativePath string    `json:"relativePath,omitempty"`
	Size         int64     `json:"size,omitempty"`
	DateAdded    time.Time `json:"dateAdded,omitempty"`
	SceneName    string    `json:"sceneName,omitempty"`
	ReleaseGroup string    `json:"releaseGroup,omitempty"`
	Quality      Quality   `json:"quality,omitempty"`
	ID           int       `json:"id,omitempty"`
	MediaInfo    MediaInfo `json:"mediaInfo,omitempty"`
}

type Quality struct {
	Quality struct {
		ID         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Source     string `json:"source,omitempty"`
		Resolution int    `json:"resolution,omitempty"`
		Modifier   string `json:"modifier,omitempty"`
	} `json:"quality,omitempty"`
	Revision struct {
		Version  int  `json:"version,omitempty"`
		Real     int  `json:"real,omitempty"`
		IsRepack bool `json:"isRepack,omitempty"`
	} `json:"revision,omitempty"`
}

type AlternativeTitle struct {
	SourceType string `json:"sourceType,omitempty"`
	MovieID    int    `json:"movieId,omitempty"`
	Title      string `json:"title,omitempty"`
	SourceID   int    `json:"sourceId,omitempty"`
	Votes      int    `json:"votes,omitempty"`
	VoteCount  int    `json:"voteCount,omitempty"`
	Language   struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"language,omitempty"`
	ID int `json:"id,omitempty"`
}

type Ratings struct {
	Votes int     `json:"votes,omitempty"`
	Value float64 `json:"value,omitempty"`
}

type MovieTag struct {
	Label string `json:"label,omitempty"`
	Id    int    `json:"id,omitempty"`
}

type Image struct {
	CoverType string `json:"coverType,omitempty"`
	URL       string `json:"url,omitempty"`
}

type Profile struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type Folder struct {
	Path      string `json:"path,omitempty"`
	FreeSpace int64  `json:"freeSpace,omitempty"`
	ID        int    `json:"id,omitempty"`
}

type AddMovieRequest struct {
	Title             string          `json:"title,omitempty"`
	TitleSlug         string          `json:"titleSlug,omitempty"`
	Images            []Image         `json:"images,omitempty"`
	QualityProfileID  int             `json:"qualityProfileId,omitempty"`
	LanguageProfileID int             `json:"languageProfileId,omitempty"`
	TMDBID            int             `json:"tmdbId,omitempty"`
	RootFolderPath    string          `json:"rootFolderPath,omitempty"`
	Monitored         bool            `json:"monitored,omitempty"`
	AddOptions        AddMovieOptions `json:"addOptions,omitempty"`
	Year              int             `json:"year,omitempty"`
	Tags              []int           `json:"tags,omitempty"`
}

type AddMovieOptions struct {
	SearchForMovie bool `json:"searchForMovie,omitempty"`
}

type SystemStatus struct {
	Version           string `json:"version,omitempty"`
	BuildTime         string `json:"buildTime,omitempty"`
	IsDebug           bool   `json:"isDebug,omitempty"`
	IsProduction      bool   `json:"isProduction,omitempty"`
	IsAdmin           bool   `json:"isAdmin,omitempty"`
	IsUserInteractive bool   `json:"isUserInteractive,omitempty"`
	StartupPath       string `json:"startupPath,omitempty"`
	AppData           string `json:"appData,omitempty"`
	OsName            string `json:"osName,omitempty"`
	OsVersion         string `json:"osVersion,omitempty"`
	IsMonoRuntime     bool   `json:"isMonoRuntime,omitempty"`
	IsMono            bool   `json:"isMono,omitempty"`
	IsLinux           bool   `json:"isLinux,omitempty"`
	IsOsx             bool   `json:"isOsx,omitempty"`
	IsWindows         bool   `json:"isWindows,omitempty"`
	Branch            string `json:"branch,omitempty"`
	Authentication    string `json:"authentication,omitempty"`
	SqliteVersion     string `json:"sqliteVersion,omitempty"`
	UrlBase           string `json:"urlBase,omitempty"`
	RuntimeVersion    string `json:"runtimeVersion,omitempty"`
	RuntimeName       string `json:"runtimeName,omitempty"`
}

type MediaInfo struct {
	AudioAdditionalFeatures      string  `json:"audioAdditionalFeatures,omitempty"`
	AudioBitrate                 int     `json:"audioBitrate,omitempty"`
	AudioChannelPositions        string  `json:"audioChannelPositions,omitempty"`
	AudioChannelPositionsText    string  `json:"audioChannelPositionsText,omitempty"`
	AudioChannels                int     `json:"audioChannels,omitempty"`
	AudioCodecID                 string  `json:"audioCodecID,omitempty"`
	AudioCodecLibrary            string  `json:"audioCodecLibrary,omitempty"`
	AudioFormat                  string  `json:"audioFormat,omitempty"`
	AudioLanguages               string  `json:"audioLanguages,omitempty"`
	AudioProfile                 string  `json:"audioProfile,omitempty"`
	AudioStreamCount             int     `json:"audioStreamCount,omitempty"`
	ContainerFormat              string  `json:"containerFormat,omitempty"`
	Height                       int     `json:"height,omitempty"`
	RunTime                      string  `json:"runTime,omitempty"`
	ScanType                     string  `json:"scanType,omitempty"`
	SchemaRevision               int     `json:"schemaRevision,omitempty"`
	Subtitles                    string  `json:"subtitles,omitempty"`
	VideoBitDepth                int     `json:"videoBitDepth,omitempty"`
	VideoBitrate                 int     `json:"videoBitrate,omitempty"`
	VideoCodecID                 string  `json:"videoCodecID,omitempty"`
	VideoCodecLibrary            string  `json:"videoCodecLibrary,omitempty"`
	VideoColourPrimaries         string  `json:"videoColourPrimaries,omitempty"`
	VideoFormat                  string  `json:"videoFormat,omitempty"`
	VideoFps                     float64 `json:"videoFps,omitempty"`
	VideoMultiViewCount          int     `json:"videoMultiViewCount,omitempty"`
	VideoProfile                 string  `json:"videoProfile,omitempty"`
	VideoTransferCharacteristics string  `json:"videoTransferCharacteristics,omitempty"`
	Width                        int     `json:"width,omitempty"`
}
