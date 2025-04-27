package listviewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
)

type Model struct {
	ctx            *context.Context
	viewport       viewport.Model
	NumRows        int
	currentId      int
	bottomBoundId  int
	topBoundId     int
	listItemHeight int
}

func NewModel(ctx *context.Context, dimensions constants.Dimensions, numRows, itemHeight int) Model {
	m := Model{
		ctx:            ctx,
		viewport:       viewport.New(dimensions.Width, dimensions.Height),
		NumRows:        numRows,
		currentId:      0,
		topBoundId:     0,
		listItemHeight: itemHeight,
	}
	m.bottomBoundId = min(m.GetNumItemsDisplayed()-1, m.NumRows-1)
	return m
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) SetNumRows(numRows int) {
	m.NumRows = numRows
	m.bottomBoundId = min(m.GetNumItemsDisplayed()-1, m.NumRows-1)
}

func (m *Model) NextItem() int {
	if m.currentId == m.bottomBoundId {
		m.viewport.ScrollDown(m.listItemHeight)
		m.topBoundId += 1
		m.bottomBoundId += 1
	}

	m.currentId = min(m.currentId+1, m.NumRows-1)
	return m.currentId
}

func (m *Model) PrevItem() int {
	if m.currentId <= m.topBoundId {
		m.viewport.ScrollUp(m.listItemHeight)
		m.topBoundId -= 1
		m.bottomBoundId -= 1
	}
	m.currentId = max(m.currentId-1, 0)
	return m.currentId
}

func (m *Model) FirstItem() int {
	m.currentId = 0
	m.viewport.GotoTop()
	return m.currentId
}

func (m *Model) LastItem() int {
	m.currentId = m.NumRows - 1
	m.viewport.GotoBottom()
	return m.currentId
}

func (m *Model) GetCurrItem() int {
	return m.currentId
}

func (m *Model) GetNumItemsDisplayed() int {
	if m.listItemHeight == 0 {
		return 0
	}
	return m.viewport.Height / m.listItemHeight
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
