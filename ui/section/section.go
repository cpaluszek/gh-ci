package section

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/components/table"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
)

type BaseModel struct {
	Title   string
	Ctx     context.Context
	Table   table.Model
	Columns []table.Column
}

type Component interface {
	Update(msg tea.Msg) (Section, tea.Cmd)
	View() string
}

type Section interface {
	Table
	Component
}

type Table interface {
	NumRows() int
	// GetCurrRow() data.RowData
	CurrRow() int
	NextRow() int
	PrevRow() int
	BuildRows() []table.Row
	// GetIsLoading() bool
	// SetIsLoading(val bool)
}

func NewModel(
	ctx context.Context,
	title string,
	columns []table.Column,
) BaseModel {
	m := BaseModel{
		Title:   title,
		Ctx:     ctx,
		Columns: columns,
	}

	m.Table = table.NewModel(
		constants.Dimensions{
			Width:  ctx.ScreenWidth,
			Height: ctx.ScreenHeight,
		},
		columns,
		nil,
	)

	return m
}

func (m *BaseModel) NextRow() int {
	return m.Table.NextItem()
}

func (m *BaseModel) PrevRow() int {
	return m.Table.PrevItem()
}

func (m *BaseModel) CurrRow() int {
	return m.Table.GetCurrItem()
}

func (m *BaseModel) View() string {
	if m.Table.Rows == nil {
		return lipgloss.Place(
			m.Ctx.ScreenWidth,
			m.Ctx.ScreenHeight,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s\n\nNo data available", m.Title),
		)
	}
	return m.Table.View()
}
