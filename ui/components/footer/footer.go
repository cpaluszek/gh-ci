package footer

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	content string
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
	return m.content
}

func (m Model) Update(msg string) (Model, tea.Cmd) {
	return m, nil
}
