package reposection

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/commands"
	"github.com/cpaluszek/gh-ci/ui/components/table"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
	"github.com/cpaluszek/gh-ci/ui/section"
)

type Model struct {
	section.BaseModel
	repos []*github.Repository
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
				Grow:  false,
			},
			{
				Title: "Stars",
				Width: 10,
				Grow:  false,
			},
			{
				Title: "Visibility",
				Width: 20,
				Grow:  false,
			},
			{
				Title: "Last Updated",
				Width: 20,
				Grow:  false,
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
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch  {
		case key.Matches(msg, keys.Keys.OpenGitHub):
			if m.repos == nil {
				return m, nil
			}
			currentIndex := m.Table.GetCurrItem()
			if currentIndex < 0 || currentIndex >= len(m.repos) {
				return m, nil
			}
			url := m.repos[currentIndex].URL
			if url == "" {
				return m, nil
			}

			return m, commands.OpenBrowser(url)
		}
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
	for _, repo := range m.repos {
		language := repo.Language
		stars := fmt.Sprintf("%d", repo.StargazerCount)
		updated := repo.UpdatedAt.Format("Jan 2, 2006")
		visibility := ""
		if repo.IsPrivate {
			visibility = "Private"
		} else {
			visibility = "Public"
		}
		rows = append(rows, table.Row{
			repo.Name,
			language,
			stars,
			visibility,
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
	fetchCmd := commands.FetchRepositories(m.Ctx.Client, m.Ctx.Config.Github.Repositories)
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
