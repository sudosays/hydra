package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
)

type PneumaUI struct {
	Screen           tcell.Screen
	cursorX, cursorY int
	Exit             bool
	Style            tcell.Style
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
		Screen:  screen,
		cursorX: 0,
		cursorY: 0,
		Exit:    false,
	}
	return ui
}

func (ui PneumaUI) Close() {
	ui.Screen.Fini()
	os.Exit(0)
	ui.Exit = true
}

func (ui PneumaUI) Tick() {
	ui.Screen.Sync()
	switch ev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyRune {
			if ev.Rune() == 'q' {
				ui.Close()
			}
		}
	}
}

func (ui *PneumaUI) MoveCursor(x, y int) error {
	w, h := ui.Screen.Size()
	if x < 0 || x > w || y < 0 || y > h {
		return errors.New("PneumaUI: cursor out of bounds.")
	}

	ui.cursorX = x
	ui.cursorY = y
	return nil
}

func (ui PneumaUI) PutRune(r rune) {
	ui.Screen.SetContent(ui.cursorX, ui.cursorY, r, []rune{}, tcell.StyleDefault)
}

func (ui PneumaUI) PutStr(str string) {
	for _, c := range str {
		ui.PutRune(rune(c))
		ui.cursorX++
	}
}

func (ui PneumaUI) PrintList(list []string) {
	for num, item := range list {
		outStr := fmt.Sprintf("%d) %s", num+1, item)
		ui.PutStr(outStr)
		ui.cursorY++
	}
}

func (ui PneumaUI) ShowFooter(content string) {
	cursorYPos := ui.cursorY
	_, h := ui.Screen.Size()
	ui.cursorY = h - 1
	ui.PutStr(content)
	ui.cursorY = cursorYPos
}
