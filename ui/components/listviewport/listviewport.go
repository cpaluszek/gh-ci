package listviewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
)

type Model struct {
	ctx       *context.Context
	viewport  viewport.Model
	NumRows   int
	currentId int
}

func NewModel(ctx *context.Context, dimensions constants.Dimensions, numRows int) Model {
	return Model{
		ctx:       ctx,
		viewport:  viewport.New(dimensions.Width, dimensions.Height),
		NumRows:   numRows,
		currentId: 0,
	}
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) NextItem() int {
	m.currentId = min(m.currentId+1, m.NumRows-1)
	return m.currentId
}

func (m *Model) PrevItem() int {
	m.currentId = max(m.currentId-1, 0)
	return m.currentId
}

func (m *Model) GetCurrItem() int {
	return m.currentId
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.viewport.Width = dimensions.Width
	m.viewport.Height = dimensions.Height
}

func (m *Model) View() string {
	viewport := m.viewport.View()
	return lipgloss.NewStyle().
		Width(m.viewport.Width).
		MaxWidth(m.viewport.Width).
		Render(viewport)
}

func (m *Model) UpdateContext(ctx *context.Context) {
	m.ctx = ctx
}
