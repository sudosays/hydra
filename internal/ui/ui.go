package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	//"strings"
	"log"
)

// Mode defines the behaviour of the UI: Navigate or Input. It controls how
// EventKeys are handled in the Tick() method.
type Mode int

// Modes for the UI are: Navigate, and Input.
const (
	Navigate Mode = iota
	Input
)

// CommandKey stores a tcell compatible key that can be compared to an EventKey
// event. The Rune for a non-alphabetical key must be set as the Rune value for
// that key. e.g. KeyEnter == rune(13)
type CommandKey struct {
	Key  tcell.Key
	Rune rune
	Mod  tcell.ModMask
}

type cursor struct {
	X, Y int
}

// PneumaUI represents a UI that the user can interact with. It contains a
// tcell.Screen as well as commands and content.
type PneumaUI struct {
	Screen      tcell.Screen
	Cursor      cursor
	Exit        bool
	Style       tcell.Style
	Mode        Mode
	InputBuffer string
	Content     []Drawable
	Commands    map[CommandKey]func()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Init creates everything necessary for an interactive user-interface by
// creating and initialising a new tcell.Screen, blank command map and setting
// the cursor to 0,0. It returns a PneumaUI strict.
func Init() PneumaUI {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	screen.Clear()
	commands := make(map[CommandKey]func())
	ui := PneumaUI{
		Screen:   screen,
		Cursor:   cursor{0, 0},
		Exit:     false,
		Mode:     Navigate,
		Style:    tcell.StyleDefault,
		Commands: commands,
	}
	return ui
}

// Redraw clears the screen and re-renders all of the widgets in content. It
// does not reset the cursor or any other state of the PneumaUI
func (ui *PneumaUI) Redraw() {
	ui.Screen.Clear()
	ui.Draw()
	ui.Screen.Sync()
}

// Reset clears the screen, content and cursor. It essentially puts the ui into
// the same state as Init does.
func (ui *PneumaUI) Reset() {
	ui.Screen.Clear()
	ui.Screen.Sync()
	ui.Cursor = cursor{0, 0}
	ui.Mode = Navigate
	ui.InputBuffer = ""
	ui.Content = make([]Drawable, 0)
}

// Close finalises the tcell.Screen and exits the program with status code 0.
func (ui PneumaUI) Close() {
	ui.Screen.Fini()
	os.Exit(0)
	ui.Exit = true
}

// WaitForInput blocks until the user finished entering a string, and returns
// the inputBuffer of the PneumaUI.
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

// SetCommands takes a map of CommandKeys and callback functions that will be
// checked against EventKeys in Tick().
func (ui *PneumaUI) SetCommands(commands map[CommandKey]func()) {
	ui.Commands = commands
}

// Tick is the main method that applications using the ui should spin on. It
// makes sure to check against the Commands while in Navigate mode, executing
// callbacks if necessary. During Input mode, it ensures to capture
// alphanumerical letters to the input buffer (terminating on escape or enter)
// and rendering the character to screen at the cursor.
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
				ui.putRune(ev.Rune())
				ui.Cursor.X++
			}
		}
	}
}

// MoveCursor attempts to place the cursor at a given location and checks for
// out of bounds errors.
func (ui *PneumaUI) MoveCursor(x, y int) error {
	w, h := ui.Screen.Size()
	if x < 0 || x > w || y < 0 || y > h {
		return errors.New("cursor out of bounds")
	}

	ui.Cursor.X = x
	ui.Cursor.Y = y
	return nil
}

func (ui PneumaUI) putRune(r rune) {
	//ui.Screen.SetContent(ui.Cursor.X, ui.Cursor.Y, r, []rune{}, tcell.StyleDefault)
	ui.Screen.SetContent(ui.Cursor.X, ui.Cursor.Y, r, []rune{}, ui.Style)
}

func (ui *PneumaUI) putString(str string) {
	for _, c := range str {
		ui.putRune(rune(c))
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

	ui.hLine(0, w, h-2)

	ui.Cursor.Y = h - 2
	ui.Cursor.X = 0
	ui.Cursor.Y = h - 1
	ui.putString(fmt.Sprintf("%-*s", w, footerContent))
	ui.Cursor.Y = cursorYPos
	ui.Cursor.X = cursorXPos
}

func (ui *PneumaUI) hLine(startX, endX, y int) {
	oldCursorPos := ui.Cursor
	for x := startX; x < endX; x++ {
		ui.MoveCursor(x, y)
		ui.putString("─")
	}
	ui.Cursor = oldCursorPos
}

func (ui *PneumaUI) vLine(startY, endY, x int) {
	oldCursorPos := ui.Cursor
	for y := startY; y < endY; y++ {
		ui.MoveCursor(x, y)
		ui.putString("│")

	}
	ui.Cursor = oldCursorPos
}

func (ui *PneumaUI) box(x, y, w, h int) {
	corners := []string{"╭", "╮", "╰", "╯"}
	ui.hLine(x, x+w, y)
	ui.hLine(x, x+w, y+h)
	ui.vLine(y, y+h, x)
	ui.vLine(y, y+h, x+w)
	ui.MoveCursor(x, y)
	ui.putString(corners[0])
	ui.MoveCursor(x+w, y)
	ui.putString(corners[1])
	ui.MoveCursor(x, y+h)
	ui.putString(corners[2])
	ui.MoveCursor(x+w, y+h)
	ui.putString(corners[3])
}

// Draw renders all of the content to the screen, in order. It does not care
// about clearing things that are already on screen, meaning that successive
// Draw() calls can be made with modified content without replacing anything on
// the screen.
func (ui *PneumaUI) Draw() {
	for _, drawable := range ui.Content {
		drawable.Draw(ui)
	}
}

// AddLabel creates a new Label widget an adds it to the content.
func (ui *PneumaUI) AddLabel(x, y int, text string) *Label {
	label := &Label{X: x, Y: y, Content: text}
	ui.Content = append(ui.Content, label)
	return label
}

// AddTable appends a new Table widget to the content.
func (ui *PneumaUI) AddTable(x, y int, headings []string, content [][]string) *Table {
	table := &Table{X: x, Y: y, Headings: headings, Content: content, Active: true, Index: 0}
	ui.Content = append(ui.Content, table)
	return table
}
