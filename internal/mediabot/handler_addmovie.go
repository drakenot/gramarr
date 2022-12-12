package mediabot

import (
	"fmt"
	"github.com/drakenot/gramarr/pkg/telegram"
	"strings"

	"github.com/drakenot/gramarr/internal/util"
	"github.com/drakenot/gramarr/pkg/chatbot"
	"github.com/drakenot/gramarr/pkg/radarr"

	"path/filepath"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) HandleAddMovie(m *tb.Message) {
	mb.StartConversation(NewAddMovieConversation(mb), m)
}

func NewAddMovieConversation(b *MediaBot) *AddMovieConversation {
	return &AddMovieConversation{bot: b, Client: b.tele}
}

type AddMovieConversation struct {
	*telegram.Client
	currentStep    chatbot.Handler
	movieQuery     string
	movieResults   []radarr.Movie
	folderResults  []radarr.Folder
	selectedMovie  *radarr.Movie
	selectedFolder *radarr.Folder
	bot            *MediaBot
}

func (c *AddMovieConversation) Run(m *tb.Message) {
	c.currentStep = c.AskMovie(m)
}

func (c *AddMovieConversation) Name() string {
	return "addmovie"
}

func (c *AddMovieConversation) CurrentStep() chatbot.Handler {
	return c.currentStep
}

func (c *AddMovieConversation) AskMovie(m *tb.Message) chatbot.Handler {
	c.Send(m.Sender, "What movie do you want to search for?")

	return func(m *tb.Message) {
		c.movieQuery = m.Text

		movies, err := c.bot.radarr.SearchMovies(c.movieQuery)
		c.movieResults = movies

		// Search Service Failed
		if err != nil {
			c.bot.tele.SendError(m.Sender, "Failed to search movies.")
			c.bot.StopConversation(c)
			return
		}

		// No Results
		if len(movies) == 0 {
			msg := fmt.Sprintf("No movie found with the title '%s'", util.EscapeMarkdown(c.movieQuery))
			c.Send(m.Sender, msg)
			c.bot.StopConversation(c)
			return
		}

		// Found some movies! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d movies:*", len(movies)))
		for i, movie := range movies {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(movie.String())))
		}
		c.Send(m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickMovie(m)
	}
}

func (c *AddMovieConversation) AskPickMovie(m *tb.Message) chatbot.Handler {

	// Send custom reply keyboard
	var options []string
	for _, movie := range c.movieResults {
		options = append(options, fmt.Sprintf("%s", movie))
	}
	options = append(options, "/cancel")
	c.SendChoices(m.Sender, "Which one would you like to download?", options)

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
			c.SendError(m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickMovie(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddMovieConversation) AskFolder(m *tb.Message) chatbot.Handler {

	user, _ := c.bot.users.User(m.Sender.ID)

	folders, err := c.bot.radarr.GetFolders(user.IsAdmin())
	c.folderResults = folders

	// GetFolders Service Failed
	if err != nil {
		c.SendError(m.Sender, "Failed to get folders.")
		c.bot.StopConversation(c)
		return nil
	}

	// No Results
	if len(folders) == 0 {
		c.SendError(m.Sender, "No destination folders found.")
		c.bot.StopConversation(c)
		return nil
	}

	// Found folders!

	// Send the results
	var msg []string
	msg = append(msg, fmt.Sprintf("*Found %d folders:*", len(folders)))
	for i, folder := range folders {
		msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(filepath.Base(folder.Path))))
	}
	c.Send(m.Sender, strings.Join(msg, "\n"))

	// Send the custom reply keyboard
	var options []string
	for _, folder := range folders {
		options = append(options, fmt.Sprintf("%s", filepath.Base(folder.Path)))
	}
	options = append(options, "/cancel")
	c.SendChoices(m.Sender, "Which folder should it download to?", options)

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
			c.SendError(m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddMovie(m)
	}
}

func (c *AddMovieConversation) AddMovie(m *tb.Message) {
	_, err := c.bot.radarr.AddMovie(*c.selectedMovie, c.bot.config.Radarr.QualityID, c.selectedFolder.Path, util.GetUserName(m))

	// Failed to add movie
	if err != nil {
		c.SendError(m.Sender, "Failed to add movie.")
		c.bot.StopConversation(c)
		return
	}

	if c.selectedMovie.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(c.selectedMovie.RemotePoster)}
		c.Send(m.Sender, photo)
	}

	// Notify User
	c.Send(m.Sender, "Movie has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added movie '%s'", util.DisplayName(m.Sender), util.EscapeMarkdown(c.selectedMovie.String()))
	c.SendAdmin(c.bot.users.Admins(), adminMsg)

	c.bot.StopConversation(c)
}
