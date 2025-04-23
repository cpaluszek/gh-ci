package sidebar

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/styles"
)

type Model struct {
	ctx *context.Context
	viewport viewport.Model
	data string
}

func NewModel(ctx *context.Context) Model {
	return Model {
		data: "",
		viewport: viewport.Model{
			Width: 0,
			Height: 0,
		},
		ctx: ctx,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	height := m.ctx.MainContentHeight
	width := constants.SideBarWidth

	style := styles.SideBarStyle.
		Height(height).
		MaxHeight(height).
		Width(width).
		MaxWidth(width)

	if m.data == "" {
		return style.Align(lipgloss.Center).Render(
			lipgloss.PlaceVertical(height, lipgloss.Center, "No data...."),
		)
	}

	return style.Render(m.viewport.View())
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent(data)
}

func (m *Model) UpdateProgramContext(ctx *context.Context) {
	if ctx == nil {
		return
	}

	m.ctx = ctx
	m.viewport.Height = m.ctx.MainContentHeight
	m.viewport.Width = constants.SideBarWidth
}


