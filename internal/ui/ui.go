package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	//"strings"
	"log"
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
	Content     []Drawable
	Commands    map[CommandKey]func()
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
	commands := make(map[CommandKey]func())
	ui := PneumaUI{
		Screen:   screen,
		Cursor:   Cursor{0, 0},
		Exit:     false,
		Mode:     Navigate,
		Style:    tcell.StyleDefault,
		Commands: commands,
	}
	return ui
}

func (ui *PneumaUI) Redraw() {
	ui.Screen.Clear()
	ui.Draw()
	ui.Screen.Sync()
}

func (ui *PneumaUI) Reset() {
	ui.Screen.Clear()
	ui.Screen.Sync()
	ui.Cursor = Cursor{0, 0}
	ui.Mode = Navigate
	ui.InputBuffer = ""
	ui.Content = make([]Drawable, 0)
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

func (ui *PneumaUI) SetCommands(commands map[CommandKey]func()) {
	ui.Commands = commands
}

func (ui *PneumaUI) Tick() {
	ui.drawFooter()
	ui.Screen.Sync()
	switch ev := ui.Screen.PollEvent().(type) {
	case *tcell.EventResize:
		ui.Redraw()
	case *tcell.EventKey:
		if ui.Mode == Navigate {
			log.SetOutput(os.Stderr)
			log.Printf("Rune of key is %v\n", ev.Rune())
			cmd := CommandKey{Key: ev.Key(), Rune: ev.Rune(), Mod: ev.Modifiers()}
			if callback, ok := ui.Commands[cmd]; ok {
				callback()
				ui.Redraw()
			}
		} else if ui.Mode == Input {
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
	//ui.Screen.SetContent(ui.Cursor.X, ui.Cursor.Y, r, []rune{}, tcell.StyleDefault)
	ui.Screen.SetContent(ui.Cursor.X, ui.Cursor.Y, r, []rune{}, ui.Style)
}

func (ui *PneumaUI) PutStr(str string) {
	for _, c := range str {
		ui.PutRune(rune(c))
		ui.Cursor.X++
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

func (ui *PneumaUI) Box(x, y, w, h int) {
	corners := []string{"╭", "╮", "╰", "╯"}
	ui.HLine(x, x+w, y)
	ui.HLine(x, x+w, y+h)
	ui.VLine(y, y+h, x)
	ui.VLine(y, y+h, x+w)
	ui.MoveCursor(x, y)
	ui.PutStr(corners[0])
	ui.MoveCursor(x+w, y)
	ui.PutStr(corners[1])
	ui.MoveCursor(x, y+h)
	ui.PutStr(corners[2])
	ui.MoveCursor(x+w, y+h)
	ui.PutStr(corners[3])
}

// Debug function to call specific
func (ui *PneumaUI) Draw() {
	for _, drawable := range ui.Content {
		drawable.Draw(ui)
	}
}

func (ui *PneumaUI) AddLabel(x, y int, text string) *Label {
	label := &Label{X: x, Y: y, Content: text}
	ui.Content = append(ui.Content, label)
	return label
}

func (ui *PneumaUI) AddTable(x, y int, headings []string, content [][]string) *Table {
	table := &Table{X: x, Y: y, Headings: headings, Content: content, Active: true, Index: 0}
	ui.Content = append(ui.Content, table)
	return table
}
