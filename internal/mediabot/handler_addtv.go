package mediabot

import (
	"fmt"
	"github.com/drakenot/gramarr/pkg/telegram"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/drakenot/gramarr/internal/util"
	"github.com/drakenot/gramarr/pkg/chatbot"
	"github.com/drakenot/gramarr/pkg/sonarr"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (mb *MediaBot) HandleAddTVShow(m *tb.Message) {
	mb.StartConversation(NewAddTVShowConversation(mb), m)
}
func NewAddTVShowConversation(b *MediaBot) *AddTVShowConversation {
	return &AddTVShowConversation{bot: b, Client: b.tele}
}

type AddTVShowConversation struct {
	*telegram.Client
	bot                     *MediaBot
	currentStep             chatbot.Handler
	TVQuery                 string
	TVShowResults           []sonarr.TVShow
	folderResults           []sonarr.Folder
	selectedTVShow          *sonarr.TVShow
	selectedTVShowSeasons   []int
	selectedQualityProfile  *sonarr.Profile
	selectedLanguageProfile *sonarr.Profile
	selectedFolder          *sonarr.Folder
}

func (c *AddTVShowConversation) Run(m *tb.Message) {
	c.currentStep = c.AskTVShow(m)
}

func (c *AddTVShowConversation) Name() string {
	return "addtv"
}

func (c *AddTVShowConversation) CurrentStep() chatbot.Handler {
	return c.currentStep
}

func (c *AddTVShowConversation) AskTVShow(m *tb.Message) chatbot.Handler {
	c.Send(m.Sender, "What TV Show do you want to search for?")

	return func(m *tb.Message) {
		c.TVQuery = m.Text

		TVShows, err := c.bot.sonarr.SearchTVShows(c.TVQuery)
		c.TVShowResults = TVShows

		// Search Service Failed
		if err != nil {
			c.SendError(m.Sender, "Failed to search TV Show.")
			c.bot.StopConversation(c)
			return
		}

		// No Results
		if len(TVShows) == 0 {
			msg := fmt.Sprintf("No TV Show found with the title '%s'", util.EscapeMarkdown(c.TVQuery))
			c.Send(m.Sender, msg)
			c.bot.StopConversation(c)
			return
		}

		// Found some TVShows! Yay!
		var msg []string
		msg = append(msg, fmt.Sprintf("*Found %d TV Shows:*", len(TVShows)))
		for i, TV := range TVShows {
			msg = append(msg, fmt.Sprintf("%d) %s", i+1, util.EscapeMarkdown(TV.String())))
		}
		c.Send(m.Sender, strings.Join(msg, "\n"))
		c.currentStep = c.AskPickTVShow(m)
	}
}

func (c *AddTVShowConversation) AskPickTVShow(m *tb.Message) chatbot.Handler {

	// Send custom reply keyboard
	var options []string
	for _, TVShow := range c.TVShowResults {
		options = append(options, fmt.Sprintf("%s", TVShow))
	}
	options = append(options, "/cancel")
	c.SendChoices(m.Sender, "Which one would you like to download?", options)

	return func(m *tb.Message) {

		// Set the selected TVShow
		for i := range options {
			if m.Text == options[i] {
				c.selectedTVShow = &c.TVShowResults[i]
				break
			}
		}

		// Not a valid TV selection
		if c.selectedTVShow == nil {
			c.SendError(m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickTVShow(m)
			return
		}

		c.currentStep = c.AskPickTVShowSeason(m)
	}
}

func (c *AddTVShowConversation) isSelectedSeason(s *sonarr.TVShowSeason) bool {

	for _, season := range c.selectedTVShowSeasons {
		if s.SeasonNumber == season {
			return true
		}
	}

	return false
}

func (c *AddTVShowConversation) AskPickTVShowQuality(m *tb.Message) chatbot.Handler {

	profiles, err := c.bot.sonarr.GetProfile("qualityprofile")

	// GetProfile Service Failed
	if err != nil {
		c.SendError(m.Sender, "Failed to get quality profiles.")
		c.bot.StopConversation(c)
		return nil
	}

	// Send custom reply keyboard
	var options []string
	for _, QualityProfile := range profiles {
		options = append(options, fmt.Sprintf("%v", QualityProfile.Name))
	}
	options = append(options, "/cancel")
	c.SendChoices(m.Sender, "Which quality shall I look for?", options)

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
			c.SendError(m.Sender, "Invalid selection.")
			c.currentStep = c.AskPickTVShowQuality(m)
			return
		}

		c.currentStep = c.AskFolder(m)
	}
}

func (c *AddTVShowConversation) AskPickTVShowSeason(m *tb.Message) chatbot.Handler {

	if c.selectedTVShowSeasons == nil {
		c.selectedTVShowSeasons = []int{}
	}

	// Send custom reply keyboard
	var options []string
	options = append(options, "All")
	if len(c.selectedTVShowSeasons) > 0 {
		options = append(options, "Nope. I'm done!")
	}

	for _, season := range c.selectedTVShow.Seasons {
		if !c.isSelectedSeason(season) && season.SeasonNumber > 0 {
			options = append(options, fmt.Sprintf("%v", season.SeasonNumber))
		}
	}

	options = append(options, "/cancel")
	if len(c.selectedTVShowSeasons) > 0 {
		c.SendChoices(m.Sender, "Any other season?", options)
	} else {
		c.SendChoices(m.Sender, "Which season would you like to download?", options)
	}

	return func(m *tb.Message) {
		if m.Text == "All" {
			c.selectedTVShowSeasons = []int{}
			for _, season := range c.selectedTVShow.Seasons {
				if season.SeasonNumber > 0 {
					c.selectedTVShowSeasons = append(c.selectedTVShowSeasons, season.SeasonNumber)
				}
			}
			c.currentStep = c.AskFolder(m)
			return
		} else if m.Text == "Nope. I'm done!" {
			c.currentStep = c.AskFolder(m)
			return
		} else {
			var selectedSeason *sonarr.TVShowSeason

			// Set the selected TV
			i, err := strconv.Atoi(m.Text)
			if err == nil && i <= len(c.selectedTVShow.Seasons) && i > 0 {
				for _, season := range c.selectedTVShow.Seasons {
					if i == season.SeasonNumber {
						selectedSeason = season
						break
					}
				}
			}

			// Not a valid TV selection
			if selectedSeason == nil {
				c.SendError(m.Sender, "Invalid selection.")
				c.currentStep = c.AskPickTVShowSeason(m)
				return
			}

			c.selectedTVShowSeasons = append(c.selectedTVShowSeasons, selectedSeason.SeasonNumber)
			c.currentStep = c.AskPickTVShowSeason(m)
		}
	}
}

func (c *AddTVShowConversation) AskFolder(m *tb.Message) chatbot.Handler {

	folders, err := c.bot.sonarr.GetFolders()
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
		if c.selectedTVShow == nil {
			c.SendError(m.Sender, "Invalid selection.")
			c.currentStep = c.AskFolder(m)
			return
		}

		c.AddTVShow(m)
	}
}

func (c *AddTVShowConversation) AddTVShow(m *tb.Message) {

	_, err := c.bot.sonarr.AddTVShow(*c.selectedTVShow, sonarr.AddSeriesOptions{
		TVDBID:         c.selectedTVShow.TvdbID,
		Title:          c.selectedTVShow.Title,
		Seasons:        c.selectedTVShowSeasons,
		SeasonFolder:   true,
		RootFolderPath: c.selectedFolder.Path,
		Monitored:      true,
		SearchNow:      true,
	})

	// Failed to add TV
	if err != nil {
		c.SendError(m.Sender, "Failed to add TV.")
		c.bot.StopConversation(c)
		return
	}

	c.selectedTVShow.RemotePoster = c.bot.sonarr.GetPosterURL(*c.selectedTVShow)
	if c.selectedTVShow.RemotePoster != "" {
		photo := &tb.Photo{File: tb.FromURL(c.selectedTVShow.RemotePoster)}
		c.bot.tele.Send(m.Sender, photo)
	}

	// Notify User
	c.Send(m.Sender, "TV Show has been added!")

	// Notify Admin
	adminMsg := fmt.Sprintf("%s added TV Show '%s'", util.DisplayName(m.Sender), util.EscapeMarkdown(c.selectedTVShow.String()))
	c.SendAdmin(c.bot.users.Admins(), adminMsg)

	c.bot.StopConversation(c)
}
