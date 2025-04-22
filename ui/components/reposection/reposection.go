package reposection

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/table"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/section"
)

type Model struct {
	section.BaseModel
	repos []*github.RepositoryData
}

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Repositories",
		[]table.Column{
			{
				Title: "Repository",
				Width: 30,
				Grow:  true,
			},
			{
				Title: "Language",
				Width: 20,
				Grow:  true,
			},
			{
				Title: "Stars",
				Width: 10,
				Grow:  true,
			},
			{
				Title: "Last Updated",
				Width: 20,
				Grow:  true,
			},
		},
	)

	return Model{
		BaseModel: base,
		repos:     nil,
	}
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.RepositoriesMsg:
		m.repos = msg.Repositories
		m.SetIsLoading(false)
		m.Table.SetRows(m.BuildRows())
	}

	table, cmd := m.Table.Update(msg)
	m.Table = table
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) BuildRows() []table.Row {
	var rows []table.Row
	for _, repoData := range m.repos {
		repo := repoData.Repository
		language := ""
		if repo.Language != nil {
			language = *repo.Language
		}
		stars := "0"
		if repo.StargazersCount != nil {
			stars = fmt.Sprintf("%d", *repo.StargazersCount)
		}
		updated := "Unknown"
		if repo.UpdatedAt != nil {
			updated = repo.UpdatedAt.Format("Jan 2, 2006")
		}
		rows = append(rows, table.Row{
			*repo.FullName,
			language,
			stars,
			updated,
		})
	}
	return rows
}

func (m *Model) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.Ctx.MainContentWidth,
		Height: m.Ctx.MainContentHeight,
	}
}

func (m *Model) NumRows() int {
	return len(m.repos)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}

func (m *Model) UpdateContext(ctx *context.Context) {
	m.Ctx = ctx
	m.Table.UpdateContext(ctx)
	m.Table.SetDimensions(m.GetDimensions())
	m.Table.SyncViewPortContent()
}

func (m *Model) Fetch() []tea.Cmd {
	if m == nil {
		return nil
	}

	var cmds []tea.Cmd
	tableCmd := m.Table.StartLoadingSpinner()
	fetchCmd := commands.FetchRepositories(m.Ctx.Client)
	cmds = append(cmds, tableCmd, fetchCmd)
	m.SetIsLoading(true)
	return cmds
}

func (m *Model) GetCurrentRow() github.RowData {
	if len(m.repos) == 0 {
		return nil
	}
	return m.repos[m.Table.GetCurrItem()]
}
