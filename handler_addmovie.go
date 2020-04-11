package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gramarr/radarr"

	"path/filepath"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (e *Env) HandleAddMovie(m *tb.Message) {
	e.CM.StartConversation(NewAddMovieConversation(e), m)
}

func NewAddMovieConversation(e *Env) *AddMovieConversation {
	return &AddMovieConversation{env: e}
}

type AddMovieConversation struct {
	currentStep            Handler
	movieQuery             string
	movieResults           []radarr.Movie
	folderResults          []radarr.Folder
	selectedMovie          *radarr.Movie
	selectedQualityProfile *radarr.Profile
	selectedFolder         *radarr.Folder
	user                   User
	env                    *Env
}

func (c *AddMovieConversation) Run(m *tb.Message) {
	c.currentStep = c.AskMovie(m)
}

func (c *AddMovieConversation) Name() string {
	return "addmovie"
}

func (c *AddMovieConversation) CurrentStep() Handler {
	return c.currentStep
}

func (c *AddMovieConversation) AskMovie(m *tb.Message) Handler {
	c.user, _ = c.env.Users.User(m.Sender.ID)

	Send(c.env.Bot, m.Sender, "What movie do you want to search for?")

	return func(m *tb.Message) {
		c.movieQuery = m.Text
		movies, err := c.env.Radarr.SearchMovies(c.movieQuery)
		c.movieResults = movies

		if err != nil {
			SendError(c.env.Bot, m.Sender, "Failed to search movies.")
			c.env.CM.StopConversation(c)
			return
		}

		if len(movies) == 0 {
			Send(c.env.Bot, m.Sender, fmt.Sprintf("No movie found with the title '%s'", EscapeMarkdown(c.movieQuery)))
			c.env.CM.StopConversation(c)
			return
		}

		msg := []string{fmt.Sprintf("*Found %d movies:*", len(movies))}
		for _, movie := range movies {
			msg = append(msg, fmt.Sprintf("- %s", EscapeMarkdown(movie.String())))
		}
		Send(c.env.Bot, m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickMovie(m)
	}
}

func (c *AddMovieConversation) AskPickMovie(m *tb.Message) Handler {

	// Send custom reply keyboard
	var options []string
	for _, movie := range c.movieResults {
		options = append(options, fmt.Sprintf("%s", movie))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Which one would you like to download?", options)

	return func(m *tb.Message) {

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
			c.currentStep = c.AskPickMovie(m)
			return
		}

		// Check if movie already exists
		var existingMovie radarr.Movie
		movies, err := c.env.Radarr.GetMovies()
		if err == nil {
			for _, movie := range movies {
				if movie.TmdbID == c.selectedMovie.TmdbID {
					existingMovie = movie
				}
			}
		}
		if existingMovie.ID > 0 {
			Send(c.env.Bot, m.Sender, "This movie has already been requested. You will be added to the requester list")
			_, _ = c.env.Radarr.AddRequester(existingMovie, m.Sender.FirstName)
			m.Payload = strconv.Itoa(c.selectedMovie.ID)
			c.env.HandleDetails(m)

			c.env.CM.StopConversation(c)
			return
		} else {
			c.currentStep = c.AskPickMovieQuality(m)
		}
	}
}

func (c *AddMovieConversation) AskPickMovieQuality(m *tb.Message) Handler {
	profiles, err := c.env.Radarr.GetProfile(c.user.IsAdmin())

	// GetProfile Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get quality profiles.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(profiles) == 0 {
		SendError(c.env.Bot, m.Sender, "No profiles found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(profiles) == 1 {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("Profile '%s' has automatically been selected", profiles[0].Name))
		c.selectedQualityProfile = &profiles[0]
		c.currentStep = c.AskFolder(m)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, profile := range profiles {
		options = append(options, fmt.Sprintf("%s", profile.Name))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Select a profile", options)

	return func(m *tb.Message) {
		// Set the selected option
		for i := range options {
			if m.Text == options[i] {
				c.selectedQualityProfile = &profiles[i]
				break
			}
		}

		// Not a valid selection
		if c.selectedQualityProfile == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickMovieQuality(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddMovieConversation) AskFolder(m *tb.Message) Handler {

	folders, err := c.env.Radarr.GetFolders(c.user.IsAdmin())
	c.folderResults = folders

	// GetFolders Service Failed
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to get folders.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(folders) == 0 {
		SendError(c.env.Bot, m.Sender, "No destination folders found.")
		c.env.CM.StopConversation(c)
		return nil
	}

	if len(folders) == 1 {
		Send(c.env.Bot, m.Sender, fmt.Sprintf("Folder '%s' has automatically been selected", strings.Title(EscapeMarkdown(filepath.Base(folders[0].Path)))))
		c.selectedFolder = &folders[0]
		c.AddMovie(m)
		return nil
	}

	// Send the results
	var options []string
	for _, folder := range folders {
		options = append(options, strings.Title(fmt.Sprintf("%s", filepath.Base(folder.Path))))
	}
	options = append(options, "/cancel")
	SendKeyboardList(c.env.Bot, m.Sender, "Select a folder", options)

	return func(m *tb.Message) {
		// Set the selected folder
		for i, opt := range options {
			if m.Text == opt {
				c.selectedFolder = &c.folderResults[i]
				break
			}
		}

		// Not a valid folder selection
		if c.selectedMovie == nil {
			SendError(c.env.Bot, m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddMovie(m)
	}
}

func (c *AddMovieConversation) AddMovie(m *tb.Message) {
	_, err := c.env.Radarr.AddMovie(*c.selectedMovie, c.selectedQualityProfile.ID, c.selectedFolder.Path, m.Sender.FirstName)

	// Failed to add movie
	if err != nil {
		SendError(c.env.Bot, m.Sender, "Failed to add movie.")
		c.env.CM.StopConversation(c)
		return
	}

	if c.selectedMovie.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(c.selectedMovie.RemotePoster)}
		c.env.Bot.Send(m.Sender, photo)
	}

	// Notify User
	Send(c.env.Bot, m.Sender, "Movie has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added movie '%s'", DisplayName(m.Sender), EscapeMarkdown(c.selectedMovie.String()))
	SendAdmin(c.env.Bot, c.env.Users.Admins(), adminMsg)

	c.env.CM.StopConversation(c)
}
