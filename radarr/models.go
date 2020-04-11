package radarr

import (
	"fmt"
	"time"
)

type Movie struct {
	Added                 time.Time          `json:"added,omitempty"`
	AlternativeTitles     []AlternativeTitle `json:"alternativeTitles,omitempty"`
	CleanTitle            string             `json:"cleanTitle,omitempty"`
	Downloaded            bool               `json:"downloaded,omitempty"`
	FolderName            string             `json:"folderName,omitempty"`
	Genres                []string           `json:"genres,omitempty"`
	HasFile               bool               `json:"hasFile,omitempty"`
	ID                    int                `json:"id,omitempty"`
	Images                []Image            `json:"images,omitempty"`
	ImdbID                string             `json:"imdbId,omitempty"`
	InCinemas             time.Time          `json:"inCinemas,omitempty"`
	IsAvailable           bool               `json:"isAvailable,omitempty"`
	LastInfoSync          time.Time          `json:"lastInfoSync,omitempty"`
	MinimumAvailability   string             `json:"minimumAvailability,omitempty"`
	Monitored             bool               `json:"monitored,omitempty"`
	MovieFile             MovieFile          `json:"movieFile,omitempty"`
	Overview              string             `json:"overview,omitempty"`
	Path                  string             `json:"path,omitempty"`
	PathState             string             `json:"pathState,omitempty"`
	PhysicalRelease       time.Time          `json:"physicalRelease,omitempty"`
	ProfileID             int                `json:"profileId,omitempty"`
	QualityProfileID      int                `json:"qualityProfileId,omitempty"`
	Ratings               Ratings            `json:"ratings,omitempty"`
	RemotePoster          string             `json:"remotePoster,omitempty"`
	Runtime               int                `json:"runtime,omitempty"`
	SecondaryYearSourceID int                `json:"secondaryYearSourceId,omitempty"`
	SizeOnDisk            int64              `json:"sizeOnDisk,omitempty"`
	SortTitle             string             `json:"sortTitle,omitempty"`
	Status                string             `json:"status,omitempty"`
	Studio                string             `json:"studio,omitempty"`
	Tags                  []int              `json:"tags,omitempty"`
	Title                 string             `json:"title,omitempty,omitempty"`
	TitleSlug             string             `json:"titleSlug,omitempty"`
	TmdbID                int                `json:"tmdbId,omitempty"`
	Website               string             `json:"website,omitempty"`
	Year                  int                `json:"year,omitempty"`
	YouTubeTrailerID      string             `json:"youTubeTrailerId,omitempty"`
}

func (m Movie) String() string {
	if m.Year != 0 {
		return fmt.Sprintf("%s (%d)", m.Title, m.Year)
	} else {
		return m.Title
	}
}

type MovieFile struct {
	DateAdded    time.Time `json:"dateAdded,omitempty"`
	Edition      string    `json:"edition"`
	ID           int       `json:"id,omitempty"`
	MediaInfo    MediaInfo `json:"mediaInfo,omitempty"`
	MovieID      int       `json:"movieId,omitempty"`
	Quality      Quality   `json:"quality,omitempty"`
	RelativePath string    `json:"relativePath,omitempty"`
	ReleaseGroup string    `json:"releaseGroup,omitempty"`
	SceneName    string    `json:"sceneName,omitempty"`
	Size         int64     `json:"size,omitempty"`
}

type Quality struct {
	Quality struct {
		ID         int    `json:"id,omitempty"`
		Modifier   string `json:"modifier,omitempty"`
		Name       string `json:"name,omitempty"`
		Resolution int    `json:"resolution,omitempty"`
		Source     string `json:"source,omitempty"`
	} `json:"quality,omitempty"`
	Revision struct {
		IsRepack bool `json:"isRepack,omitempty"`
		Real     int  `json:"real,omitempty"`
		Version  int  `json:"version,omitempty"`
	} `json:"revision,omitempty"`
}

type AlternativeTitle struct {
	ID       int `json:"id,omitempty"`
	Language struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"language,omitempty"`
	MovieID    int    `json:"movieId,omitempty"`
	SourceID   int    `json:"sourceId,omitempty"`
	SourceType string `json:"sourceType,omitempty"`
	Title      string `json:"title,omitempty"`
	VoteCount  int    `json:"voteCount,omitempty"`
	Votes      int    `json:"votes,omitempty"`
}

type Ratings struct {
	Value float64 `json:"value,omitempty"`
	Votes int     `json:"votes,omitempty"`
}

type MovieTag struct {
	Id    int    `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
}

type Image struct {
	CoverType string `json:"coverType,omitempty"`
	URL       string `json:"url,omitempty"`
}

type Profile struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Folder struct {
	ID        int    `json:"id,omitempty"`
	FreeSpace int64  `json:"freeSpace,omitempty"`
	Path      string `json:"path,omitempty"`
}

type AddMovieRequest struct {
	AddOptions        AddMovieOptions `json:"addOptions,omitempty"`
	Images            []Image         `json:"images,omitempty"`
	LanguageProfileID int             `json:"languageProfileId,omitempty"`
	Monitored         bool            `json:"monitored,omitempty"`
	QualityProfileID  int             `json:"qualityProfileId,omitempty"`
	RootFolderPath    string          `json:"rootFolderPath,omitempty"`
	TMDBID            int             `json:"tmdbId,omitempty"`
	Tags              []int           `json:"tags,omitempty"`
	Title             string          `json:"title,omitempty"`
	TitleSlug         string          `json:"titleSlug,omitempty"`
	Year              int             `json:"year,omitempty"`
}

type AddMovieOptions struct {
	SearchForMovie bool `json:"searchForMovie,omitempty"`
}

type SystemStatus struct {
	AppData           string `json:"appData,omitempty"`
	Authentication    string `json:"authentication,omitempty"`
	Branch            string `json:"branch,omitempty"`
	BuildTime         string `json:"buildTime,omitempty"`
	IsAdmin           bool   `json:"isAdmin,omitempty"`
	IsDebug           bool   `json:"isDebug,omitempty"`
	IsLinux           bool   `json:"isLinux,omitempty"`
	IsMono            bool   `json:"isMono,omitempty"`
	IsMonoRuntime     bool   `json:"isMonoRuntime,omitempty"`
	IsOsx             bool   `json:"isOsx,omitempty"`
	IsProduction      bool   `json:"isProduction,omitempty"`
	IsUserInteractive bool   `json:"isUserInteractive,omitempty"`
	IsWindows         bool   `json:"isWindows,omitempty"`
	OsName            string `json:"osName,omitempty"`
	OsVersion         string `json:"osVersion,omitempty"`
	RuntimeName       string `json:"runtimeName,omitempty"`
	RuntimeVersion    string `json:"runtimeVersion,omitempty"`
	SqliteVersion     string `json:"sqliteVersion,omitempty"`
	StartupPath       string `json:"startupPath,omitempty"`
	UrlBase           string `json:"urlBase,omitempty"`
	Version           string `json:"version,omitempty"`
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
