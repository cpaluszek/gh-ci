package footer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	content string
	width   int
}

func NewModel() Model {
	return Model{
		content: " ↑/↓: navigate · enter: select · q: quit",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Render(m.content)
}

func (m Model) Update(msg string) (Model, tea.Cmd) {
	return m, nil
}

func (m *Model) SetWidth(width int) {
	m.width = width
}
