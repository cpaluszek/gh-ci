package footer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/ui/styles"
)

type Model struct {
	content              string
	width                int
	ShowQuitConfirmation bool
	quitConfirmation     string
}

func NewModel() Model {
	return Model{
		content:              " ↑/↓: navigate · enter: select · o: open · q: quit",
		ShowQuitConfirmation: false,
		quitConfirmation:     "Press q/esc again to quit",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.ShowQuitConfirmation {
		return styles.StatusBarStyle.
			Width(m.width).
			Render(m.quitConfirmation)
	}
	return styles.StatusBarStyle.
		Width(m.width).
		Render(m.content)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.ShowQuitConfirmation {
				return m, tea.Quit
			} else {
				m.ShowQuitConfirmation = true
			}
		default:
			m.ShowQuitConfirmation = false
		}
	}
	return m, nil
}

func (m *Model) SetWidth(width int) {
	m.width = width
}
