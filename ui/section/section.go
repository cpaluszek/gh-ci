package section

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/gh-actions/github"
	"github.com/cpaluszek/gh-actions/ui/components/table"
	"github.com/cpaluszek/gh-actions/ui/constants"
	"github.com/cpaluszek/gh-actions/ui/context"
)

type BaseModel struct {
	Title     string
	Ctx       *context.Context
	Table     table.Model
	Columns   []table.Column
	IsLoading bool
}

type Component interface {
	Update(msg tea.Msg) (Section, tea.Cmd)
	View() string
}

type Section interface {
	Table
	Component
	UpdateContext(ctx *context.Context)
}

type Table interface {
	NumRows() int
	GetCurrentRow() github.RowData
	CurrRow() int
	NextRow() int
	PrevRow() int
	BuildRows() []table.Row
	GetIsLoading() bool
	SetIsLoading(val bool)
	// TODO: if not all section implement this, remove it
	Fetch() []tea.Cmd
}

func NewModel(
	ctx *context.Context,
	title string,
	columns []table.Column,
) BaseModel {
	m := BaseModel{
		Title:   title,
		Ctx:     ctx,
		Columns: columns,
	}

	m.Table = table.NewModel(
		ctx,
		constants.Dimensions{
			Width:  ctx.MainContentHeight,
			Height: ctx.MainContentWidth,
		},
		columns,
		nil,
		false,
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

func (m *BaseModel) GetIsLoading() bool {
	return m.IsLoading
}

func (m *BaseModel) View() string {
	return m.Ctx.Styles.SectionContainer.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			m.Table.View(),
		),
	)
}

func (m *BaseModel) UpdateContext(ctx *context.Context) {
	m.Ctx = ctx
	m.Table.SetDimensions(
		constants.Dimensions{
			Width:  ctx.MainContentWidth,
			Height: ctx.MainContentHeight,
		},
	)
	m.Table.SyncViewPortContent()
}

func (m *BaseModel) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.Ctx.ScreenWidth,
		Height: m.Ctx.ScreenHeight,
	}
}
