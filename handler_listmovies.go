package main

import (
	"fmt"
	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	_ "github.com/fbiville/markdown-table-formatter/pkg/markdown"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

// HandleAddMovie func
func (e *Env) HandleListMovies(m *tb.Message) {
	movies, err := e.Radarr.GetMovies()
	if err != nil {
		Send(e.Bot, m.Sender, "Error loading movies")
		return
	}

	var msg []string
	if len(movies) > 0 {

		var tableArr [][]string
		for _, movie := range movies {
			tableArr = append(tableArr, []string{EscapeMarkdown(movie.Title), strconv.Itoa(movie.Year)})
		}

		prettyPrintedTable, err := markdown.NewTableFormatterBuilder().
			WithPrettyPrint().
			Build("Title", "Year").
			Format(tableArr)
		if err != nil {
			// ... do your thing
		}

		msg = append(msg, "```")
		msg = append(msg, fmt.Sprintln(prettyPrintedTable))
		msg = append(msg, "```")

	} else {
		msg = append(msg, "No movies found")
	}
	Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
}
