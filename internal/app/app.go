package app

import (
	"github.com/cpaluszek/pipeye/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width int
	height int
}

func New(cfg *config.Config) *Model {
	return &Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return "got width"
}

func (m *Model) Run() error {
	p := tea.NewProgram(
		*m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		)

	_, err := p.Run()
	return err
}
