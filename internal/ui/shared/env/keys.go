package env

import (
	"slices"

	tea "charm.land/bubbletea/v2"
)

const (
	KeyEnter     = "enter"
	KeyEsc       = "esc"
	KeyBackspace = "backspace"
	KeyCtrlC     = "ctrl+c"
	KeyTab       = "tab"
	KeyShiftTab  = "shift+tab"
	KeyCtrlT     = "ctrl+t"
)

var backKeys = map[Mode][]string{
	NormalMode: {"b", KeyBackspace, KeyEsc},
	InputMode:  {KeyEsc},
	ScanMode:   {},
}

var quitKeys = map[Mode][]string{
	NormalMode: {"q", KeyCtrlC},
	InputMode:  {KeyCtrlC},
	ScanMode:   {KeyCtrlC},
}

func IsBackKey(msg tea.KeyPressMsg, mode Mode) bool {
	if keys, ok := backKeys[mode]; ok {
		return slices.Contains(keys, msg.String())
	}
	return false
}

func IsQuitKey(msg tea.KeyPressMsg, mode Mode) bool {
	if keys, ok := quitKeys[mode]; ok {
		return slices.Contains(keys, msg.String())
	}
	return false
}
