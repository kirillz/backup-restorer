package model

import (
	"fmt"

	"backup-restorer/pkg/backup"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	programNameStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("202")).Bold(true)
)

const (
	cursorModeText  = "Enter the path to the backup directory"
	quitKeySequence = "ctrl+c, q"
	programName     = "PostgreSQL Backup Restorer"
)

type Model struct {
	textInput  textinput.Model
	progress   progress.Model
	cursorMode string
	quitting   bool
	backupDir  string
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = cursorModeText
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		textInput:  ti,
		progress:   p,
		cursorMode: cursorModeText,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if m.cursorMode == cursorModeText {
				m.backupDir = m.textInput.Value()
				if m.backupDir != "" {
					return m, tea.Batch(
						backup.RestoreBackupCmd(m.backupDir),
						m.progress.IncrPercent(10), // Увеличиваем прогресс на 10%
					)
				}
			}
		}
	case tea.WindowSizeMsg:
		m.textInput.Width = msg.Width
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case backup.RestoreBackupMsg:
		m.progress.SetPercent(1.0) // Устанавливаем прогресс на 100%
		m.quitting = true
		return m, tea.Quit
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		programNameStyle.Render(programName),
		m.textInput.View(),
		m.progress.View(),
		helpStyle.Render("Press "+quitKeySequence+" to quit."),
	)
}
