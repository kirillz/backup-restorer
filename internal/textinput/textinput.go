package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model = textinput.Model

func New() Model {
	return textinput.New()
}

var Blink = textinput.Blink
