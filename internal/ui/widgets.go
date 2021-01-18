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

	ui.Box(t.X, t.Y, maxWidth+1, len(t.Content)+1)

	ui.Style = ui.Style.Bold(true)
	ui.Style = ui.Style.Underline(true)

	ui.MoveCursor(t.X+1, t.Y+1)
	for i, heading := range t.Headings {
		ui.PutStr(fmt.Sprintf("%-*s", colWidths[i], heading))
	}

	ui.Style = ui.Style.Bold(false)
	ui.Style = ui.Style.Underline(false)

	ui.MoveCursor(t.X+1, t.Y+2)
	for i, row := range t.Content {
		for col, item := range row {

			ui.PutStr(fmt.Sprintf("%-*s", colWidths[col], item))
		}
		ui.MoveCursor(t.X+1, t.Y+2+i)
	}

	ui.Style = ui.Style.Normal()

}
