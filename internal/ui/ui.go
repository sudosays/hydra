package ui

import (
	// "fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sudosays/hydra/pkg/data/hugo"

	"strconv"
)

type model struct {
	table table.Model
}

func InitialModel(posts []hugo.Post) model {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Title", Width: 25},
		{Title: "Date", Width: 12},
		{Title: "Status", Width: 15},
	}

	rows := []table.Row{}

	for i, post := range posts {
		status := ""
		if post.Draft {
			status = "draft"
		} else {
			status = "published"
		}
		rows = append(rows, table.Row{strconv.Itoa((i + 1)), post.Title, post.Date, status})
	}

	table := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return model{
		table,
	}
}

func (m model) Init() tea.Cmd {

	return nil

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		}

	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd

}

func (m model) View() string {
	return m.table.View() + "\n"
}
