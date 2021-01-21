package ui

import (
	"github.com/gdamore/tcell/v2"
)

type CommandKey struct {
	Key  tcell.Key
	Rune rune
	Mod  tcell.ModMask
}
