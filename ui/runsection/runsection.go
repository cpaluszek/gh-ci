package runsection

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/commands"
	"github.com/cpaluszek/gh-ci/ui/components/table"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/section"
	"github.com/cpaluszek/gh-ci/ui/utils"
)

type Model struct {
	section.BaseModel
	Runs *github.WorkflowRun
}

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Run",
		[]table.Column{
			{
				Title: "Job",
				Width: 20,
				Grow:  true,
			},
			{
				Title: "Status",
				Width: 18,
				Grow:  false,
			},
			{
				Title: "Duration",
				Width: 12,
				Grow:  false,
			},
		})

	return Model{
		BaseModel: base,
		Runs:      nil,
	}
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.WorkflowRunMsg:
		m.Runs = msg.RunWithJobs
		m.Table.SetRows(m.BuildRows())
		m.Table.FirstItem()
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch msg.String() {
		case "o":
			if m.Runs == nil {
				return m, nil
			}
			currentIndex := m.Table.GetCurrItem()
			if currentIndex < 0 || currentIndex >= len(m.Runs.Jobs) {
				return m, nil
			}
			url := m.Runs.Jobs[currentIndex].URL
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
	if m.Runs == nil {
		return nil
	}

	var rows []table.Row
	for _, job := range m.Runs.Jobs {
		status := utils.GetJobStatusSymbol(m.Ctx, job.Status, job.Conclusion) + " " + job.Conclusion
		status = utils.CleanANSIEscapes(status)
		rows = append(rows, table.Row{
			job.GetName(),
			m.Ctx.Styles.Default.Render(status),
			m.Ctx.Styles.Default.Render(utils.GetJobDuration(job)),
		})
	}

	return rows
}

func (m *Model) NumRows() int {
	return len(m.Runs.Jobs)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}

func (m *Model) Fetch() []tea.Cmd {
	if m == nil {
		return nil
	}
	return nil
}

func (m *Model) GetCurrentRow() github.RowData {
	if m == nil || m.Runs == nil || len(m.Runs.Jobs) == 0 {
		return nil
	}
	currentIndex := m.Table.GetCurrItem()
	if currentIndex < 0 || currentIndex >= len(m.Runs.Jobs) {
		return nil
	}
	return m.Runs.Jobs[currentIndex]
}
