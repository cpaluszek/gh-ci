package footer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/ui/context"
)

type Model struct {
	ctx                  *context.Context
	content              string
	width                int
	ShowQuitConfirmation bool
	quitConfirmation     string
}

func NewModel(ctx *context.Context) Model {
	return Model{
		ctx:                  ctx,
		content:              " ↑/↓: navigate · enter: select · o: open · q: quit",
		width:                0,
		ShowQuitConfirmation: false,
		quitConfirmation:     "Press q/esc again to quit",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	if m.ShowQuitConfirmation {
		return m.ctx.Styles.StatusBar.
			Width(m.width).
			Render(m.quitConfirmation)
	}
	return m.ctx.Styles.StatusBar.
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
