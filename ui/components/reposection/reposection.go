package reposection

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func NewModel(ctx context.Context) Model {
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.RepositoriesMsg:
		// TODO: is loading false
		m.repos = msg.Repositories
		m.Table.SetRows(m.BuildRows())

	}

	return m, cmd
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

func (m *Model) View() string {
	if m.Table.Rows == nil {
		d := m.GetDimensions()
		return lipgloss.Place(
			d.Width,
			d.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s\n\nNo data available", m.Title),
		)
	}
	return m.Table.View()
}

func (m *Model) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.Ctx.ScreenWidth,
		Height: m.Ctx.ScreenHeight,
	}
}

func (m *Model) NumRows() int {
	return len(m.repos)
}
