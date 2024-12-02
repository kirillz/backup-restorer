package progress

import (
	"github.com/charmbracelet/bubbles/progress"
)

type Model = progress.Model

func New(opts ...progress.Option) Model {
	return progress.New(opts...)
}
