package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sudosays/hydra/internal/types"
	"github.com/sudosays/hydra/pkg/data/hugo"
	// "log"
	// "os"
	"os/exec"
	"strconv"
)

type model struct {
	editor          types.EditorCommand
	blog            hugo.Blog
	table           table.Model
	input           textinput.Model
	err             error
	editing         bool
	altScreenActive bool
	focusTable      bool
	draftsOnly      bool
}

type editorFinishedMsg struct{ err error }

func InitialModel(blog hugo.Blog, editorCmd types.EditorCommand) model {

	draftsOnly := true

	posts := []hugo.Post{}

	if draftsOnly {
		posts = blog.DraftPosts()
	} else {
		posts = blog.Posts
	}

	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Title", Width: 25},
		{Title: "Date", Width: 12},
		{Title: "Status", Width: 15},
	}

	rows := formatTableRows(posts)

	textInput := textinput.New()
	textInput.Width = 20

	table := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	return model{
		editorCmd,
		blog,
		table,
		textInput,
		nil,
		false,
		false,
		true,
		true,
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

		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.focusTable {
				return m, tea.Quit
			}

		case "esc":
			if !m.focusTable {
				m.focusTable = true
				m.table.Focus()
				m.input.Blur()
				m.input.Reset()
			}

		case "enter":
			if m.focusTable {
				path := m.blog.Posts[m.table.Cursor()-1].Path
				m.altScreenActive = true
				m.editing = true
				return m, tea.Batch(tea.ClearScreen, tea.EnterAltScreen, m.startEditor(path))
			} else {
				m.focusTable = true
				m.table.Focus()
				m.input.Blur()
			}

		case "n":
			m.focusTable = false
			m.input.Focus()
			m.table.Blur()
			return m, nil

		case "v":
			if m.draftsOnly {
				m.table.SetRows(formatTableRows(m.blog.Posts))
				m.draftsOnly = false
			} else {
				m.table.SetRows(formatTableRows(m.blog.DraftPosts()))
				m.draftsOnly = true
			}

		}

	case editorFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.altScreenActive = false
		return m, tea.Batch(tea.ClearScreen, tea.ExitAltScreen)
	}

	if m.focusTable {
		m.table, cmd = m.table.Update(msg)
	} else {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd

}

func (m model) View() string {

	var inputStyle lipgloss.Style

	if !m.focusTable {
		inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("228")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Width(52)
	} else {
		inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666")).
			Foreground(lipgloss.Color("#666")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true).
			Width(52)
	}

	statusStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Foreground(lipgloss.Color("#666")).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true).
		Width(80).
		MarginBottom(1).
		Align(lipgloss.Center)

	return lipgloss.JoinVertical(0,
		statusStyle.Render("HYDRA v0.2"),
		m.table.View(),
		"Commands: [q]uit [n]ew [e]dit cycle [v]iew\n",
		inputStyle.Render(m.input.View()))

}

func (m model) startEditor(path string) tea.Cmd {
	editorCmd := exec.Command(m.editor.Command, m.editor.Args, path)
	// editorCmd.Stdin = os.Stdin
	// editorCmd.Stdout = os.Stdout
	return tea.ExecProcess(editorCmd, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})

}

func formatTableRows(posts []hugo.Post) []table.Row {

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

	return rows

}
