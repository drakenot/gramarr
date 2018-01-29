package radarr

import (
	"fmt"
	"strings"
	"regexp"
	"net/url"

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
	}

	return u.String()
}

type Client struct {
	apiKey   string
	username string
	password string
	baseURL string
	maxResults int
	client *resty.Client
}


func (c *Client) SearchMovies(term string) ([]Movie, error) {
	resp, err := c.client.R().SetQueryParam("term", term).SetResult([]Movie{}).Get("movies/lookup")
	if err != nil {
		return nil, err
	}

	movies := *resp.Result().(*[]Movie)
	if len(movies) > c.maxResults {
		movies = movies[:c.maxResults]
	}
	return movies, nil
}

func (c *Client) GetFolders() ([]Folder, error) {
	resp, err := c.client.R().SetResult([]Folder{}).Get("rootfolder")
	if err != nil {
		return nil, err
	}

	folders := *resp.Result().(*[]Folder)
	return folders, nil
}

func (c *Client) AddMovie(m Movie, qualityProfile int, path string) (movie Movie, err error){

	request := AddMovieRequest{
		Title: m.Title,
		TitleSlug: m.TitleSlug,
		Images: m.Images,
		QualityProfileID: qualityProfile,
		TMDBID: m.TMDBID,
		RootFolderPath: path,
		Monitored: true,
		AddOptions: AddMovieOptions{SearchForMovie:true},
	}

	resp, err := c.client.R().SetBody(request).SetResult(Movie{}).Post("movie")
	if err != nil {
		return
	}

	movie = *resp.Result().(*Movie)
	return
}