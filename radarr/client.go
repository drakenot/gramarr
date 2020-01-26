package radarr

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/resty.v1"
)

var (
	apiRgx = regexp.MustCompile(`[a-z0-9]{32}`)
)

func NewClient(c Config) (*Client, error) {
	if c.Hostname == "" {
		return nil, fmt.Errorf("hostname is empty")
	}

	if match := apiRgx.MatchString(c.APIKey); !match {
		return nil, fmt.Errorf("api key is invalid format: %s", c.APIKey)
	}

	baseURL := createApiURL(c)

	r := resty.New()
	r.SetHostURL(baseURL)
	r.SetHeader("Accept", "application/json")
	r.SetQueryParam("apikey", c.APIKey)
	if c.Username != "" && c.Password != "" {
		r.SetBasicAuth(c.Username, c.Password)
	}

	client := Client{
		apiKey:     c.APIKey,
		maxResults: c.MaxResults,
		username:   c.Username,
		password:   c.Password,
		baseURL:    baseURL,
		client:     r,
	}
	return &client, nil
}

func createApiURL(c Config) string {
	c.Hostname = strings.TrimPrefix(c.Hostname, "http://")
	c.Hostname = strings.TrimPrefix(c.Hostname, "https://")
	c.URLBase = strings.TrimPrefix(c.URLBase, "/")

	u := url.URL{}
	if c.SSL {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	if c.Port == 80 {
		u.Host = c.Hostname
	} else {
		u.Host = fmt.Sprintf("%s:%d", c.Hostname, c.Port)
	}

	if c.URLBase != "" {
		u.Path = fmt.Sprintf("%s/api", c.URLBase)
	} else {
		u.Path = "/api"
	}

	return u.String()
}

type Client struct {
	apiKey     string
	username   string
	password   string
	baseURL    string
	maxResults int
	client     *resty.Client
}

func (c *Client) SearchMovie(tmdbId int) (Movie, error) {
	var movie Movie
	resp, err := c.client.R().SetQueryParam("tmdbId", strconv.Itoa(tmdbId)).SetResult(Movie{}).Get("movie/lookup/tmdb")
	if err != nil {
		return movie, err
	}
	movie = *resp.Result().(*Movie)
	return movie, nil
}

func (c *Client) SearchMovies(term string) ([]Movie, error) {
	resp, err := c.client.R().SetQueryParam("term", term).SetResult([]Movie{}).Get("movie/lookup")
	if err != nil {
		return nil, err
	}

	movies := *resp.Result().(*[]Movie)
	if len(movies) > c.maxResults {
		movies = movies[:c.maxResults]
	}
	return movies, nil
}

func (c *Client) GetProfile(prfl string) ([]Profile, error) {

	resp, err := c.client.R().SetResult([]Profile{}).Get(prfl)
	if err != nil {
		return nil, err
	}
	profile := *resp.Result().(*[]Profile)

	return profile, nil

}

func (c *Client) GetMoviesFromFolder(folder Folder) ([]Movie, error) {
	movies, err := c.GetMovies()
	if err != nil {
		return nil, err
	}
	var ret []Movie
	for _, movie := range movies {
		if strings.HasPrefix(movie.Path, folder.Path) {
			ret = append(ret, movie)
		}
	}

	return ret, nil
}

func (c *Client) GetMovies() ([]Movie, error) {
	resp, err := c.client.R().SetResult([]Movie{}).Get("movie")
	if err != nil {
		return nil, err
	}
	movies := *resp.Result().(*[]Movie)
	return movies, nil
}

func (c *Client) GetMovie(movieId int) (Movie, error) {
	var movie Movie

	resp, err := c.client.R().SetResult(Movie{}).Get("/movie/" + strconv.Itoa(movieId))
	if err != nil {
		return movie, err
	}
	movie = *resp.Result().(*Movie)

	return movie, nil
}

func (c *Client) GetFolders() ([]Folder, error) {
	resp, err := c.client.R().SetResult([]Folder{}).Get("rootfolder")
	if err != nil {
		return nil, err
	}

	folders := *resp.Result().(*[]Folder)
	return folders, nil
}

func (c *Client) AddMovie(m Movie, qualityProfile int, path string) (movie Movie, err error) {

	request := AddMovieRequest{
		Title:            m.Title,
		TitleSlug:        m.TitleSlug,
		Images:           m.Images,
		QualityProfileID: qualityProfile,
		TMDBID:           m.TMDBID,
		RootFolderPath:   path,
		Monitored:        true,
		Year:             m.Year,
		AddOptions:       AddMovieOptions{SearchForMovie: true},
	}

	resp, err := c.client.R().SetBody(request).SetResult(Movie{}).Post("movie")
	if err != nil {
		return
	}

	movie = *resp.Result().(*Movie)
	return
}

func (c *Client) GetSystemStatus() (SystemStatus, error) {
	var systemStatus SystemStatus

	resp, err := c.client.R().SetResult(SystemStatus{}).Get("/system/status")
	if err != nil {
		return systemStatus, err
	}
	systemStatus = *resp.Result().(*SystemStatus)

	return systemStatus, nil
}

func (c *Client) GetPosterURL(movie Movie) string {
	for _, image := range movie.Images {
		if image.CoverType == "poster" {
			return image.URL
		}
	}
	return ""
}

func (c *Client) GetTags() ([]MovieTag, error) {
	resp, err := c.client.R().SetResult([]MovieTag{}).Get("tag")
	if err != nil {
		return nil, err
	}
	tags := *resp.Result().(*[]MovieTag)
	return tags, nil
}

func (c *Client) GetRequester(movie Movie) string {
	tags, err := c.GetTags()
	if err != nil {
		return ""
	}
	var requester []string
	for _, movieTagId := range movie.Tags {
		for _, tag := range tags {
			if movieTagId == tag.Id {
				requester = append(requester, strings.Title(tag.Label))
			}
		}
	}
	return strings.Join(requester, ", ")
}
