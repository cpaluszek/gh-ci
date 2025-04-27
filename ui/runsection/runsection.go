package runsection

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/table"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/section"
	"github.com/cpaluszek/pipeye/ui/utils"
)

type Model struct {
	section.BaseModel
	RunWithJobs *github.WorkflowRun
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
		BaseModel:   base,
		RunWithJobs: nil,
	}
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.WorkflowRunMsg:
		m.RunWithJobs = msg.RunWithJobs
		m.Table.SetRows(m.BuildRows())
		m.Table.FirstItem()
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch msg.String() {
		case "o":
			if m.RunWithJobs == nil {
				return m, nil
			}
			currentIndex := m.Table.GetCurrItem()
			if currentIndex < 0 || currentIndex >= len(m.RunWithJobs.Jobs) {
				return m, nil
			}
			url := m.RunWithJobs.Jobs[currentIndex].GetHTMLURL()
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
	if m.RunWithJobs == nil {
		return nil
	}

	var rows []table.Row
	for _, job := range m.RunWithJobs.Jobs {
		status := utils.GetJobStatusSymbol(m.Ctx, job.GetStatus(), job.GetConclusion()) + " " + job.GetConclusion()
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
	return len(m.RunWithJobs.Jobs)
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
	if m == nil || m.RunWithJobs == nil || len(m.RunWithJobs.Jobs) == 0 {
		return nil
	}
	currentIndex := m.Table.GetCurrItem()
	if currentIndex < 0 || currentIndex >= len(m.RunWithJobs.Jobs) {
		return nil
	}
	return &github.Job{
		Job: m.RunWithJobs.Jobs[currentIndex],
	}
}
