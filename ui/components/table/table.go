package table

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/components/listviewport"
	"github.com/cpaluszek/pipeye/ui/constants"
)

type Model struct {
	Rows         []Row
	Columns      []Column
	Dimensions   constants.Dimensions
	rowsViewport listviewport.Model
	// Empty state
	// Loading
}

type Row []string

type Column struct {
	Title string
	Width int
	Grow  bool
}

func NewModel(
	dimensions constants.Dimensions,
	columns []Column,
	rows []Row,
) Model {
	return Model{
		Rows:         rows,
		Columns:      columns,
		Dimensions:   dimensions,
		rowsViewport: listviewport.NewModel(dimensions, len(rows)),
	}
}

func (m Model) Update() (Model, tea.Cmd) {
	// TODO: loading
	return m, nil
}

func (m Model) View() string {
	header := m.renderHeader()
	body := m.renderBody()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
	)
}

func (m *Model) PrevItem() int {
	currItem := m.rowsViewport.PrevItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) NextItem() int {
	currItem := m.rowsViewport.NextItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) GetCurrItem() int {
	return m.rowsViewport.GetCurrItem()
}

func (m *Model) SetRows(rows []Row) {
	m.Rows = rows
	m.rowsViewport.NumRows = len(rows)
	m.SyncViewPortContent()
}

func (m *Model) SyncViewPortContent() {
	headerColumns := m.renderHeaderColumns()
	renderedRows := make([]string, 0, len(m.Rows))
	for i := range m.Rows {
		renderedRows = append(renderedRows, m.renderRow(i, headerColumns))
	}

	m.rowsViewport.SyncViewPort(
		lipgloss.NewStyle().
			Background(lipgloss.Color("3")).
			Render(lipgloss.JoinVertical(lipgloss.Left, renderedRows...)),
	)
}

func (m Model) renderHeaderColumns() []string {
	renderedColumns := make([]string, len(m.Columns))
	takenWidth := 0
	numGrowing := 0
	for i := range m.Columns {
		if m.Columns[i].Grow {
			numGrowing++
			continue
		}

		if m.Columns[i].Width > 0 {

			renderedColumns[i] = lipgloss.NewStyle().
				Width(m.Columns[i].Width).
				MaxWidth(m.Columns[i].Width).
				Render(m.Columns[i].Title)
			takenWidth += m.Columns[i].Width
			continue
		}

		cell := lipgloss.NewStyle().Render(m.Columns[i].Title)
		renderedColumns[i] = cell
		takenWidth += lipgloss.Width(cell)
	}

	if numGrowing == 0 {
		return renderedColumns
	}

	remainingWidth := m.Dimensions.Width - takenWidth
	growWidth := remainingWidth / numGrowing
	for i := range m.Columns {
		if !m.Columns[i].Grow {
			continue
		}
		renderedColumns[i] = lipgloss.NewStyle().
			Width(growWidth).
			MaxWidth(growWidth).
			Render(m.Columns[i].Title)
	}
	return renderedColumns
}

func (m Model) renderHeader() string {
	headerColumns := m.renderHeaderColumns()
	header := lipgloss.JoinHorizontal(
		lipgloss.Left,
		headerColumns...,
	)

	// TODO: table style
	return lipgloss.NewStyle().
		Width(m.Dimensions.Width).
		MaxWidth(m.Dimensions.Width).
		Height(m.Dimensions.Height).
		MaxHeight(m.Dimensions.Height).
		Render(header)
}

func (m Model) renderBody() string {
	bodyStyle := lipgloss.NewStyle().
		Height(m.Dimensions.Height).
		MaxWidth(m.Dimensions.Width)

	// TODO: if is loading

	if len(m.Rows) == 0 {
		return bodyStyle.Render("No data")
	}
	return m.rowsViewport.View()
}

func (m *Model) renderRow(rowId int, headerColumns []string) string {
	var style lipgloss.Style

	if m.rowsViewport.GetCurrItem() == rowId {
		style = lipgloss.NewStyle().Background(lipgloss.Color("#444444"))
	} else {
		style = lipgloss.NewStyle()
	}

	renderedColumns := make([]string, 0, len(m.Columns))
	headerColId := 0
	for i := range m.Columns {
		colWidth := lipgloss.Width(headerColumns[headerColId])
		colHeight := 1

		col := m.Rows[rowId][i]
		renderedCol := style.
			Width(colWidth).
			MaxWidth(colWidth).
			Height(colHeight).
			MaxHeight(colHeight).
			Render(col)

		renderedColumns = append(renderedColumns, renderedCol)
		headerColId++
	}

	return lipgloss.NewStyle().
		MaxWidth(m.Dimensions.Width).
		Render(
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				renderedColumns...,
			),
		)
}
