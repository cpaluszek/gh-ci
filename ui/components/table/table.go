package table

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/components/listviewport"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
)

type Model struct {
	ctx          *context.Context
	Rows         []Row
	Columns      []Column
	Dimensions   constants.Dimensions
	rowsViewport listviewport.Model
	// Empty state
	isLoading      bool
	loadingSpinner spinner.Model
}

type Row []string

type Column struct {
	Title string
	Width int
	Grow  bool
}

func NewModel(
	ctx *context.Context,
	dimensions constants.Dimensions,
	columns []Column,
	rows []Row,
	isLoading bool,
) Model {
	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.MiniDot
	loadingSpinner.Style = ctx.Styles.Spinner
	return Model{
		ctx:        ctx,
		Rows:       rows,
		Columns:    columns,
		Dimensions: dimensions,
		rowsViewport: listviewport.NewModel(
			ctx,
			constants.Dimensions{
				Width:  dimensions.Width,
				Height: dimensions.Height - constants.TableHeaderHeight,
			},
			len(rows),
			constants.TableRowHeight,
		),
		isLoading:      isLoading,
		loadingSpinner: loadingSpinner,
	}
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.isLoading {
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
	}
	return *m, cmd
}

func (m Model) StartLoadingSpinner() tea.Cmd {
	return m.loadingSpinner.Tick
}

func (m *Model) SetIsLoading(val bool) {
	m.isLoading = val
}

func (m Model) IsLoading() bool {
	return m.isLoading
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.Dimensions = dimensions
	m.rowsViewport.SetDimensions(
		constants.Dimensions{
			Width:  dimensions.Width,
			Height: dimensions.Height - constants.TableHeaderHeight,
		},
	)
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

func (m *Model) FirstItem() int {
	currItem := m.rowsViewport.FirstItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) LastItem() int {
	currItem := m.rowsViewport.LastItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) GetCurrItem() int {
	return m.rowsViewport.GetCurrItem()
}

func (m *Model) SetRows(rows []Row) {
	m.Rows = rows
	m.rowsViewport.SetNumRows(len(rows))
	m.SyncViewPortContent()
}

func (m *Model) SyncViewPortContent() {
	headerColumns := m.renderHeaderColumns()
	renderedRows := make([]string, 0, len(m.Rows))
	for i := range m.Rows {
		renderedRows = append(renderedRows, m.renderRow(i, headerColumns))
	}

	m.rowsViewport.SyncViewPort(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
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

			renderedColumns[i] = m.ctx.Styles.Header.
				Width(m.Columns[i].Width).
				MaxWidth(m.Columns[i].Width).
				Render(m.Columns[i].Title)
			takenWidth += m.Columns[i].Width
			continue
		}

		cell := m.ctx.Styles.Header.Render(m.Columns[i].Title)
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
		renderedColumns[i] = m.ctx.Styles.Header.
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

	return m.ctx.Styles.Header.
		Width(m.Dimensions.Width).
		MaxWidth(m.Dimensions.Width).
		Height(constants.TableHeaderHeight).
		MaxHeight(constants.TableHeaderHeight).
		Render(header)
}

func (m Model) renderBody() string {
	if m.isLoading {
		return lipgloss.Place(
			m.Dimensions.Width,
			m.Dimensions.Height-constants.TableHeaderHeight,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading...", m.loadingSpinner.View()),
		)
	}

	if len(m.Rows) == 0 {
		return lipgloss.Place(
			m.Dimensions.Width,
			m.Dimensions.Height-constants.TableHeaderHeight,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading...", m.loadingSpinner.View()),
		)
	}
	return m.rowsViewport.View()
}

func (m *Model) renderRow(rowId int, headerColumns []string) string {
	var style lipgloss.Style

	if m.rowsViewport.GetCurrItem() == rowId {
		style = m.ctx.Styles.SelectedRow
	} else {
		style = m.ctx.Styles.Row
	}

	renderedColumns := make([]string, 0, len(m.Columns))
	headerColId := 0
	for i := range m.Columns {
		colWidth := lipgloss.Width(headerColumns[headerColId])

		col := m.Rows[rowId][i]
		renderedCol := style.
			Width(colWidth).
			MaxWidth(colWidth).
			Height(constants.TableRowHeight).
			MaxHeight(constants.TableRowHeight).
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

func (m *Model) UpdateContext(ctx *context.Context) {
	m.ctx = ctx
	m.rowsViewport.UpdateContext(ctx)
}
