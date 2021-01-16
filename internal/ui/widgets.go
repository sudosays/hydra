package ui

import (
	"fmt"
)

type Drawable interface {
	Draw(ui *PneumaUI)
}

type Label struct {
	X, Y    int
	Content string
}

type Table struct {
	X, Y     int
	Headings []string
	Content  [][]string
	Active   bool
	Index    int
}

func (l Label) Draw(ui *PneumaUI) {
	ui.MoveCursor(l.X, l.Y)
	ui.PutStr(l.Content)
}

func (t *Table) AddRow(row []string) {
	t.Content = append(t.Content, row)
}

func (t Table) Draw(ui *PneumaUI) {

	sep := "│"
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

	ui.MoveCursor(t.X, t.Y)
	ui.PutStr(sep)
	for i, heading := range t.Headings {
		ui.PutStr(fmt.Sprintf("%-*s%s", colWidths[i], heading, sep))
	}
	ui.MoveCursor(t.X, t.Y+1)
	for i, row := range t.Content {
		ui.PutStr(sep)
		for col, item := range row {
			ui.PutStr(fmt.Sprintf("%-*s%s", colWidths[col], item, sep))
		}
		ui.MoveCursor(t.X, t.Y+1+i)
	}

}
