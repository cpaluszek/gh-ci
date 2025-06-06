package workflowssection

import (
	"sort"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/commands"
	"github.com/cpaluszek/gh-ci/ui/components/table"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
	"github.com/cpaluszek/gh-ci/ui/section"
	"github.com/cpaluszek/gh-ci/ui/utils"
)

type WorkflowRunInfo struct {
	Workflow *github.Workflow
	Run      *github.WorkflowRun
}

type Model struct {
	section.BaseModel
	workflows *github.Repository
	allRuns   []WorkflowRunInfo
}

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Workflows",
		[]table.Column{
			{
				Title: "",
				Width: 4,
				Grow:  false,
			},
			{
				Title: "Workflow",
				Width: 20,
				Grow:  false,
			},
			{
				Title: "Status",
				Width: 18,
				Grow:  false,
			},
			{
				Title: "Branch",
				Width: 32,
				Grow:  false,
			},
			{
				Title: "Triggered",
				Width: 14,
				Grow:  false,
			},
			{
				Title: "Duration",
				Width: 12,
				Grow:  false,
			},
			{
				Title: "Jobs",
				Width: 20,
				Grow:  false,
			},
			{
				Title: "Commit",
				Width: 20,
				Grow:  true,
			},
		},
	)

	return Model{
		BaseModel: base,
	}
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case commands.WorkflowsMsg:
		m.workflows = msg.Workflows
		m.allRuns = m.buildRunsList()
		m.Table.SetRows(m.BuildRows())
		m.Table.FirstItem()
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.OpenGitHub):
			if m.workflows == nil || len(m.allRuns) == 0 {
				return m, nil
			}
			currentIndex := m.Table.GetCurrItem()
			if currentIndex < 0 || currentIndex >= len(m.allRuns) {
				return m, nil
			}
			url := m.allRuns[currentIndex].Run.URL
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

func (m *Model) buildRunsList() []WorkflowRunInfo {
	var runs []WorkflowRunInfo
	for _, workflow := range m.workflows.Workflows {
		if len(workflow.Runs) == 0 {
			continue
		}
		for _, runWithJob := range workflow.Runs {
			runs = append(runs, WorkflowRunInfo{
				Workflow: workflow,
				Run:      runWithJob,
			})
		}
	}

	// Sort by creation time (most recent first)
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Run.CreatedAt.After(runs[j].Run.CreatedAt)
	})

	return runs
}

func (m Model) BuildRows() []table.Row {
	var rows []table.Row
	for _, runInfo := range m.allRuns {
		run := runInfo.Run
		workflow := runInfo.Workflow

		duration := utils.GetWorkflowRunDuration(run)
		commitMsg := run.HeadCommit.Message
		displayStatus := utils.GetWorkflowRunStatus(m.Ctx, run)

		// Build jobs indicators with symbols
		jobs := ""
		for _, job := range runInfo.Run.Jobs {
			jobs += utils.GetJobStatusSymbol(m.Ctx, job.Status, job.Conclusion)
		}
		jobs = utils.CleanANSIEscapes(jobs)

		// Table row
		rows = append(rows, table.Row{
			" " + utils.GetRunEventSymbol(m.Ctx, run.Event),
			workflow.Name,
			displayStatus,
			run.HeadBranch,
			utils.FormatTime(run.CreatedAt),
			duration,
			jobs,
			commitMsg,
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
	return len(m.allRuns)
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

	return nil
}

func (m *Model) GetCurrentRow() github.RowData {
	if m.workflows == nil || len(m.allRuns) == 0 {
		return nil
	}

	currentIndex := m.Table.GetCurrItem()
	if currentIndex < 0 || currentIndex >= len(m.allRuns) {
		return nil
	}

	return m.allRuns[currentIndex].Run
}
