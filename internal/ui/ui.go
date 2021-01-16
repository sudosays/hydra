package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	//"strings"
)

type UIMode int

const (
	Navigate UIMode = iota
	Input
)

type Cursor struct {
	X, Y int
}

type PneumaUI struct {
	Screen      tcell.Screen
	Cursor      Cursor
	Exit        bool
	Style       tcell.Style
	Mode        UIMode
	InputBuffer string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Init() PneumaUI {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	screen.Clear()
	ui := PneumaUI{
		Screen: screen,
		Cursor: Cursor{0, 0},
		Exit:   false,
		Mode:   Navigate,
	}
	return ui
}

func (ui *PneumaUI) Reset() {
	ui.Screen.Clear()
	ui.Screen.Sync()
	ui.Cursor = Cursor{0, 0}
	ui.Mode = Navigate
	ui.InputBuffer = ""
}

func (ui PneumaUI) Close() {
	ui.Screen.Fini()
	os.Exit(0)
	ui.Exit = true
}

func (ui *PneumaUI) WaitForInput() string {
	ui.Mode = Input
	for {
		ui.Tick()
		if ui.Mode == Navigate {
			break
		}
	}
	return ui.InputBuffer
}

func (ui *PneumaUI) WaitForCommand() string {
	cmd := "nop"
	return cmd
}

func (ui *PneumaUI) Tick() {
	ui.drawFooter()
	ui.Screen.Sync()
	switch ev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		if ui.Mode == Navigate {
			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'q':
					ui.Close()
				case 'i':
					ui.Mode = Input
				}
			}
		} else {
			if ev.Key() == tcell.KeyEnter || ev.Key() == tcell.KeyEscape {
				ui.Mode = Navigate
			} else if ev.Key() == tcell.KeyRune {
				ui.InputBuffer += string(ev.Rune())
				ui.PutRune(ev.Rune())
				ui.Cursor.X++
			}
		}
	}
}

func (ui *PneumaUI) MoveCursor(x, y int) error {
	w, h := ui.Screen.Size()
	if x < 0 || x > w || y < 0 || y > h {
		return errors.New("PneumaUI: cursor out of bounds.")
	}

	ui.Cursor.X = x
	ui.Cursor.Y = y
	return nil
}

func (ui PneumaUI) PutRune(r rune) {
	ui.Screen.SetContent(ui.Cursor.X, ui.Cursor.Y, r, []rune{}, tcell.StyleDefault)
}

func (ui *PneumaUI) PutStr(str string) {
	for _, c := range str {
		ui.PutRune(rune(c))
		ui.Cursor.X++
	}
}

func (ui PneumaUI) PrintList(list []string) {
	for num, item := range list {
		outStr := fmt.Sprintf("%2d) %s", num+1, item)
		ui.PutStr(outStr)
		ui.Cursor.Y++
	}
}

func (ui PneumaUI) drawFooter() {
	var footerContent string
	switch ui.Mode {
	case Navigate:
		footerContent = "Q: quit, I: insert"
	case Input:
		footerContent = "INPUT"
	}
	cursorYPos := ui.Cursor.Y
	cursorXPos := ui.Cursor.X
	w, h := ui.Screen.Size()

	ui.HLine(0, w, h-2)

	ui.Cursor.Y = h - 2
	ui.Cursor.X = 0
	ui.Cursor.Y = h - 1
	ui.PutStr(fmt.Sprintf("%-*s", w, footerContent))
	ui.Cursor.Y = cursorYPos
	ui.Cursor.X = cursorXPos
}

func (ui *PneumaUI) HLine(startX, endX, y int) {
	oldCursorPos := ui.Cursor
	for x := startX; x < endX; x++ {
		ui.MoveCursor(x, y)
		ui.PutStr("─")
	}
	ui.Cursor = oldCursorPos
}

func (ui *PneumaUI) VLine(startY, endY, x int) {
	oldCursorPos := ui.Cursor
	for y := startY; y < endY; y++ {
		ui.MoveCursor(x, y)
		ui.PutStr("│")

	}
	ui.Cursor = oldCursorPos
}

// Debug function to call specific
func (ui *PneumaUI) Draw(drawable Drawable) {
	drawable.Draw(ui)
}

func (ui *PneumaUI) AddTable(x, y int, content [][]string) *Table {
	headings := content[0]
	rows := content[0:]
	return &Table{X: x, Y: y, Headings: headings, Content: rows}
}
