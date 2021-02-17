package ui

import (
	"fmt"
)

// A Drawable is any struct (typically) that can be drawn to a PneumaUI. It
// must implement a relevant draw call by directly modifying the PneumaUI.
type Drawable interface {
	Draw(ui *PneumaUI)
}

// A Label is any text string. Essentially just a wrapper around
// PenumaUI.putStr, but with coordinates.
type Label struct {
	X, Y    int
	Content string
}

// A Table is a structured widget that contains string data in rows and
// columns. It does not manage headers or support non-string content, merely
// renders it. The Index represents the row (0 indexed) that is currently
// highlighted.
type Table struct {
	X, Y     int
	Headings []string
	Content  [][]string
	Active   bool
	Index    int
}

// Draw renders a label to the given PneumaUI.
func (l Label) Draw(ui *PneumaUI) {
	ui.MoveCursor(l.X, l.Y)
	ui.putString(l.Content)
}

// AddRow appends a row to the content of a table. It does not redraw the
// table.
func (t *Table) AddRow(row []string) {
	t.Content = append(t.Content, row)
}

// SetContent allows the table headings and rows to be changed dynamically.
// This does not trigger a redraw, but the changes will be seen upon the Draw()
// function being called.
func (t *Table) SetContent(headings []string, content [][]string) {
    t.Headings = headings
    t.Content = content
}

// Draw renders a table to the given PneumaUI. It makes sure to size the
// columns to the max width of the widest item and does not truncate the
// contents.  Furthermore, the selected item is higlighted with
// tcell.Style.Reverse
func (t Table) Draw(ui *PneumaUI) {
	if !t.Active {
		ui.Style = ui.Style.Dim(true)
	}
	colWidths := make([]int, len(t.Headings))
	for i, h := range t.Headings {
		colWidths[i] = len(h) + 1
	}
	for _, row := range t.Content {
		for col, item := range row {
			if colWidths[col] < len(item)+1 {
				colWidths[col] = len(item) + 1
			}
		}
	}

	maxWidth := 0
	for _, width := range colWidths {
		maxWidth += width
	}

	ui.box(t.X, t.Y, maxWidth+1, len(t.Content)+2)

	ui.Style = ui.Style.Bold(true)
	ui.Style = ui.Style.Underline(true)

	ui.MoveCursor(t.X+1, t.Y+1)
	for i, heading := range t.Headings {
		ui.putString(fmt.Sprintf("%-*s", colWidths[i], heading))
	}

	ui.Style = ui.Style.Bold(false)
	ui.Style = ui.Style.Underline(false)

	ui.MoveCursor(t.X+1, t.Y+2)
	for i, row := range t.Content {
		for col, item := range row {
			if i == t.Index {
				ui.Style = ui.Style.Reverse(true)
				ui.putString(fmt.Sprintf("%-*s", colWidths[col], item))
				ui.Style = ui.Style.Reverse(false)
			} else {
				ui.putString(fmt.Sprintf("%-*s", colWidths[col], item))
			}
		}
		ui.MoveCursor(t.X+1, t.Y+3+i)
	}

	ui.Style = ui.Style.Normal()

}

// NextItem sets the index of the selected item to the next one in the content
// list, wrapping around when it reaches the end.
func (t *Table) NextItem() {
	t.Index = (t.Index + 1) % len(t.Content)
}

// PreviousItem sets the index to the previous item in the content list wrapping
// around to the end when it reaches the start.
func (t *Table) PreviousItem() {
	if t.Index-1 < 0 {
		t.Index = len(t.Content) - 1
	} else {
		t.Index--
	}
}
