package main

import (
	"fmt"
	"github.com/gramarr/radarr"
	tb "gopkg.in/tucnak/telebot.v2"
	"sort"
	"strconv"
	"strings"
)

func (e *Env) HandleListMovies(m *tb.Message) {
	e.CM.StartConversation(NewListMoviesConversation(e), m)
}

func NewListMoviesConversation(e *Env) *ListMoviesConversation {
	return &ListMoviesConversation{env: e}
}

type ListMoviesConversation struct {
	currentStep            Handler
	movieQuery             string
	movieResults           []radarr.Movie
	folderResults          []radarr.Folder
	selectedMovie          *radarr.Movie
	selectedQualityProfile *radarr.Profile
	selectedFolder         *radarr.Folder
	env                    *Env
}

func (c *ListMoviesConversation) Run(m *tb.Message) {
	c.currentStep = c.AskRequester(m)
}

func (c *ListMoviesConversation) Name() string {
	return "listMovies"
}

func (c *ListMoviesConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *ListMoviesConversation) AskRequester(m *tb.Message) Handler {

	requesterList, err := c.env.Radarr.GetTags()

	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get requester list.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(requesterList) == 0 {
		SendError(c.env.Bot, m.Sender, "No requester found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	var options []string
	options = append(options, "All")
	for _, requester := range requesterList {
		options = append(options, fmt.Sprintf("%s", strings.TrimSpace(requester.Label)))
	}

	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which movies would you like to list?", options)

	return func(m *tb.Message) {
		// Set the selected folder
		for _, opt := range options {
			if m.Text == opt {
				if m.Text == "All" {
					c.movieResults, _ = c.env.Radarr.GetMovies()
				} else {
					c.movieResults, err = c.env.Radarr.GetMoviesByRequester(strings.TrimSpace(m.Text))
				}
			}
		}

		if err != nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskRequester(m)
			return
		}

		c.currentStep = c.AskMovie(m)
	}
}

func (c *ListMoviesConversation) AskMovie(m *tb.Message) Handler {

	sort.Slice(c.movieResults, func(i, j int) bool {
		return c.movieResults[i].Title < c.movieResults[j].Title
	})

	var fulfilled = []string{"*Available Movies:*"}
	var pending = []string{"*Pending Movies:*"}
	var options []string
	for _, movie := range c.movieResults {
		options = append(options, EscapeMarkdown(movie.Title))
		if movie.HasFile {
			fulfilled = append(fulfilled, fmt.Sprintf("- %s", EscapeMarkdown(movie.Title)))
		} else {
			pending = append(pending, fmt.Sprintf("- %s", EscapeMarkdown(movie.Title)))
		}
	}
	if len(fulfilled) > 1 {
		Send(c.env.Bot, m.Sender, strings.Join(fulfilled, "\n"))
	}
	if len(pending) > 1 {
		Send(c.env.Bot, m.Sender, strings.Join(pending, "\n"))
	}

	if len(options) > 0 {
		options = append(options, "Back to requester selection")
		options = append(options, "/cancel")
		SendKeyboardList(c.env.Bot, m.Sender, "Select a movie for more details", options)
	} else {
		Send(c.env.Bot, m.Sender, "No movies found")
		c.currentStep = c.AskRequester(m)
	}

	return func(m *tb.Message) {
		if m.Text == "Back to requester selection" {
			c.currentStep = c.AskRequester(m)
			return
		}

		// Set the selected movie
		for i, opt := range options {
			if m.Text == opt {
				c.selectedMovie = &c.movieResults[i]
				break
			}
		}

		// Not a valid movie selection
		if c.selectedMovie == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			return
		}

		m.Payload = strconv.Itoa(c.selectedMovie.ID)
		c.env.HandleMovieDetails(m)
		c.env.CM.StopConversation(c)
	}

}
