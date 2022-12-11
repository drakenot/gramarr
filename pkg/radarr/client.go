package radarr

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
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
	r.SetBaseURL(baseURL)
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
		u.Path = fmt.Sprintf("%s/api/v3", c.URLBase)
	} else {
		u.Path = "/api/v3"
	}

	return u.String()
}

type Client struct {
	apiKey             string
	username           string
	password           string
	baseURL            string
	maxResults         int
	restrictedFolders  []int
	restrictedProfiles []int
	client             *resty.Client
}

func (c *Client) DeleteMovie(movieId int) (err error) {
	_, err = c.client.R().SetQueryParam("deleteFiles", "true").Delete("movie/" + strconv.Itoa(movieId))
	return
}

func (c *Client) RemoveRequester(m Movie, requester string) (Movie, error) {
	tag, err := c.GetTagByLabel(requester, true)
	if err != nil {
		return m, err
	}
	var filteredTags []int
	for i := range m.Tags {
		if m.Tags[i] != tag.Id {
			filteredTags = append(filteredTags, m.Tags[i])
		}
	}
	m.Tags = filteredTags
	return c.UpdateMovie(m)
}

func (c *Client) AddRequester(m Movie, requester string) (Movie, error) {
	tag, err := c.GetTagByLabel(requester, true)
	if err != nil {
		return m, err
	}
	m.Tags = append(m.Tags, tag.Id)
	return c.UpdateMovie(m)
}

func (c *Client) UpdateMovie(m Movie) (movie Movie, err error) {
	resp, err := c.client.R().SetBody(m).SetResult(Movie{}).Put("movie")
	if err != nil {
		return
	}
	movie = *resp.Result().(*Movie)
	return
}

func (c *Client) SearchMovie(tmdbId int) (movie Movie, err error) {
	resp, err := c.client.R().SetQueryParam("tmdbId", strconv.Itoa(tmdbId)).SetResult(Movie{}).Get("movie/lookup/tmdb")
	if err != nil {
		return
	}
	movie = *resp.Result().(*Movie)
	return
}

func (c *Client) SearchMovies(term string) (movies []Movie, err error) {
	resp, err := c.client.R().SetQueryParam("term", term).SetResult([]Movie{}).Get("movie/lookup")
	if err != nil {
		return
	}
	movies = *resp.Result().(*[]Movie)
	if len(movies) > c.maxResults {
		movies = movies[:c.maxResults]
	}
	return
}

func (c *Client) GetProfile(isAdmin bool) (profiles []Profile, err error) {
	resp, err := c.client.R().SetResult([]Profile{}).Get("qualityProfile")
	if err != nil {
		return
	}
	allProfiles := *resp.Result().(*[]Profile)
	if isAdmin {
		return allProfiles, err
	}

	for _, profile := range allProfiles {
		if !contains(c.restrictedProfiles, profile.ID) {
			profiles = append(profiles, profile)
		}
	}
	return
}

func (c *Client) GetMoviesByRequester(requester string) (movies []Movie, err error) {
	allMovies, err := c.GetMovies()
	if err != nil {
		return
	}
	for _, movie := range allMovies {
		for _, t := range movie.Tags {
			tag, _ := c.GetTagById(t)
			if strings.Trim(requester, " ") == strings.Trim(tag.Label, " ") {
				movies = append(movies, movie)
			}
		}
	}
	return
}

func (c *Client) GetMoviesByFolder(folder Folder) (movies []Movie, err error) {
	allMovies, err := c.GetMovies()
	if err != nil {
		return
	}
	for _, movie := range allMovies {
		if strings.HasPrefix(movie.Path, folder.Path) {
			movies = append(movies, movie)
		}
	}
	return
}

func (c *Client) GetMovies() (movies []Movie, err error) {
	resp, err := c.client.R().SetResult([]Movie{}).Get("movie")
	if err != nil {
		return
	}
	allMovies := *resp.Result().(*[]Movie)
	for _, movie := range allMovies {
		if movie.Monitored {
			movies = append(movies, movie)
		}
	}
	return
}

func (c *Client) GetMovie(movieId int) (movie Movie, err error) {
	resp, err := c.client.R().SetResult(Movie{}).Get("movie/" + strconv.Itoa(movieId))
	if err != nil {
		return
	}
	movie = *resp.Result().(*Movie)
	return
}

func (c *Client) GetFolders(isAdmin bool) (folders []Folder, err error) {
	resp, err := c.client.R().SetResult([]Folder{}).Get("rootfolder")
	if err != nil {
		return nil, err
	}
	allFolders := *resp.Result().(*[]Folder)
	if isAdmin {
		return allFolders, nil
	}

	for _, folder := range allFolders {
		if !contains(c.restrictedFolders, folder.ID) {
			folders = append(folders, folder)
		}
	}
	return folders, nil
}

func (c *Client) AddMovie(m Movie, qualityProfile int, path string, requester string) (movie Movie, err error) {
	request := AddMovieRequest{
		Title:            m.Title,
		TitleSlug:        m.TitleSlug,
		Images:           m.Images,
		QualityProfileID: qualityProfile,
		TMDBID:           m.TmdbID,
		RootFolderPath:   path,
		Monitored:        true,
		Year:             m.Year,
		AddOptions:       AddMovieOptions{SearchForMovie: true},
	}

	tag, err := c.GetTagByLabel(requester, true)
	if err == nil {
		request.Tags = []int{tag.Id}
	}

	resp, err := c.client.R().SetBody(request).SetResult(Movie{}).Post("movie")
	if err != nil {
		return
	}

	movie = *resp.Result().(*Movie)
	return
}

func (c *Client) GetSystemStatus() (systemStatus SystemStatus, err error) {
	resp, err := c.client.R().SetResult(SystemStatus{}).Get("/system/status")
	if err != nil {
		return
	}
	systemStatus = *resp.Result().(*SystemStatus)
	return
}

func (c *Client) GetPosterURL(movie Movie) string {
	for _, image := range movie.Images {
		if image.CoverType == "poster" {
			return image.URL
		}
	}
	return ""
}

func (c *Client) GetTagByLabel(label string, createNew bool) (movieTag MovieTag, err error) {
	tags, err := c.GetTags()
	if err != nil {
		return
	}
	for _, tag := range tags {
		if strings.EqualFold(strings.TrimSpace(label), strings.TrimSpace(tag.Label)) {
			movieTag = tag
		}
	}
	if createNew && movieTag.Id == 0 {
		movieTag, err = c.CreateTag(strings.TrimSpace(label))
	}
	return
}

func (c *Client) GetTagById(id int) (movieTag MovieTag, err error) {
	tags, err := c.GetTags()
	if err != nil {
		return
	}
	for _, tag := range tags {
		if id == tag.Id {
			movieTag = tag
		}
	}
	return
}

func (c *Client) GetTags() (tags []MovieTag, err error) {
	resp, err := c.client.R().SetResult([]MovieTag{}).Get("tag")
	if err != nil {
		return
	}
	tags = *resp.Result().(*[]MovieTag)
	return
}

func (c *Client) CreateTag(label string) (tag MovieTag, err error) {
	label = strings.TrimSpace(label)
	resp, err := c.client.R().SetBody(MovieTag{Label: label}).SetResult(MovieTag{}).Post("tag")
	if err != nil {
		return
	}
	tag = *resp.Result().(*MovieTag)
	return
}

func (c *Client) GetRequesterList(movie Movie) (requester []string) {
	for _, tagId := range movie.Tags {
		tag, err := c.GetTagById(tagId)
		if err == nil {
			requester = append(requester, strings.Title(tag.Label))
		}
	}
	return
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
