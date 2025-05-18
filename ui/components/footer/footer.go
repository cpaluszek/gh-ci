package footer

import (
	bbhelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
)

type Model struct {
	ctx                  *context.Context
	content              string
	width                int
	ShowQuitConfirmation bool
	quitConfirmation     string
	Help                 bbhelp.Model
}

func NewModel(ctx *context.Context) Model {
	help := bbhelp.New()
	help.ShowAll = false
    help.Styles = ctx.Styles.Help

	return Model{
		ctx:                  ctx,
		content:              "",
		width:                0,
		ShowQuitConfirmation: false,
		quitConfirmation:     "Press q/esc again to quit",
		Help:                 help,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Quit):
			if m.ShowQuitConfirmation {
				return m, tea.Quit
			} else {
				m.ShowQuitConfirmation = true
			}
		case m.ShowQuitConfirmation && !key.Matches(msg, keys.Keys.Quit):
			m.ShowQuitConfirmation = false
		case key.Matches(msg, keys.Keys.Help):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.ShowQuitConfirmation {
		return m.ctx.Styles.Footer.Width(m.width).Render(m.quitConfirmation)
	} else {
		return m.ctx.Styles.Footer.Width(m.width).Render(m.Help.View(keys.Keys))
	}
}

func (m *Model) SetWidth(width int) {
	m.width = width
    m.Help.Width = width
}
