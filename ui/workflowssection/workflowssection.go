package workflowssection

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/table"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/section"
	"github.com/cpaluszek/pipeye/ui/styles"
	"github.com/cpaluszek/pipeye/ui/utils"
	gh "github.com/google/go-github/v71/github"
)

type WorkflowRunInfo struct {
	Workflow *gh.Workflow
	Run      *github.WorkflowRunWithJobs
}

type Model struct {
	section.BaseModel
	workflows *github.RepositoryData
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
		switch msg.String() {
		case "o":
			if m.workflows == nil || len(m.allRuns) == 0 {
				return m, nil
			}
			currentIndex := m.Table.GetCurrItem()
			if currentIndex < 0 || currentIndex >= len(m.allRuns) {
				return m, nil
			}
			url := m.allRuns[currentIndex].Run.Run.GetHTMLURL()
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
	for _, workflow := range m.workflows.WorkflowRunWithJobs {
		if len(workflow.Runs) == 0 {
			continue
		}
		for _, runWithJob := range workflow.Runs {
			runs = append(runs, WorkflowRunInfo{
				Workflow: workflow.Workflow,
				Run:      runWithJob,
			})
		}
	}

	// Sort by creation time (most recent first)
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Run.Run.GetCreatedAt().After(runs[j].Run.Run.GetCreatedAt().Time)
	})

	return runs
}

func (m Model) BuildRows() []table.Row {
	var rows []table.Row
	for _, runInfo := range m.allRuns {
		run := runInfo.Run.Run
		workflow := runInfo.Workflow

		duration := utils.GetWorkflowRunDuration(run)
		commitMsg := run.GetHeadCommit().GetMessage()
		displayStatus := utils.GetWorkflowRunStatus(run)

		// Build jobs indicators with symbols
		jobs := ""
		for _, job := range runInfo.Run.Jobs {
			jobs += styles.GetJobStatusSymbol(job.GetConclusion())
		}
		jobs = utils.CleanANSIEscapes(jobs)

		// Table row
		rows = append(rows, table.Row{
			" " + utils.GetRunEventIcon(*run.Event),
			workflow.GetName(),
			displayStatus,
			run.GetHeadBranch(),
			utils.FormatTime(run.GetCreatedAt().Time),
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
