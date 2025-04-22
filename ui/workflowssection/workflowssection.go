package workflowssection

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
	"github.com/cpaluszek/pipeye/ui/styles"
	"github.com/cpaluszek/pipeye/ui/utils"
)

type Model struct {
	section.BaseModel
	workflows *github.RepositoryData
}

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Workflows",
		[]table.Column{
			{
				Title: "Status",
				Width: 20,
				Grow:  false,
			},
			{
				Title: "Branch",
				Width: 30,
				Grow:  false,
			},
			{
				Title: "Triggered",
				Width: 30,
				Grow:  false,
			},
			{
				Title: "Duration",
				Width: 20,
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
		m.Table.SetRows(m.BuildRows())
		m.Table.FirstItem()
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
	selectedWorkflow := m.workflows.WorkflowRunWithJobs[0]
	for _, runWithJob := range selectedWorkflow.Runs {
		run := runWithJob.Run

		// Calculate duration
		var duration string
		if run.GetUpdatedAt().After(run.GetCreatedAt().Time) {
			durationTime := run.GetUpdatedAt().Sub(run.GetCreatedAt().Time)
			if durationTime.Hours() >= 1 {
				duration = fmt.Sprintf("%.1fh", durationTime.Hours())
			} else if durationTime.Minutes() >= 1 {
				duration = fmt.Sprintf("%.1fm", durationTime.Minutes())
			} else {
				duration = fmt.Sprintf("%.1fs", durationTime.Seconds())
			}
		} else {
			duration = "running"
		}

		commitMsg := run.GetHeadCommit().GetMessage()

		// Style based on conclusion
		var statusStyle lipgloss.Style
		status := run.GetStatus()
		conclusion := run.GetConclusion()

		switch status {
		case "completed":
			switch conclusion {
			case "success":
				statusStyle = styles.SuccessStyle
			case "failure", "timed_out", "startup_failure":
				statusStyle = styles.FailureStyle
			case "cancelled", "skipped", "neutral":
				statusStyle = styles.CanceledStyle
			default:
				statusStyle = styles.InProgressStyle
			}
		case "in_progress":
			statusStyle = styles.InProgressStyle
			status = "running"
		}

		displayStatus := status
		if conclusion != "" && status == "completed" {
			displayStatus = conclusion
		}

		var jobStyle lipgloss.Style
		jobs := ""
		for _, job := range runWithJob.Jobs {
			switch job.GetConclusion() {
			case "success":
				jobStyle = styles.SuccessStyle
			case "failure":
				jobStyle = styles.FailureStyle
			case "cancelled":
				jobStyle = styles.CanceledStyle
			}

			jobs += jobStyle.Render("●")
		}

		// Table row
		rows = append(rows, table.Row{
			statusStyle.Render(displayStatus),
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
	if m.workflows == nil || len(m.workflows.WorkflowRunWithJobs) == 0 {
		return 0
	}
	return len(m.workflows.WorkflowRunWithJobs[0].Runs)
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
	if m.workflows == nil || len(m.workflows.WorkflowRunWithJobs) == 0 {
		return nil
	}

	selectedWorkflow := m.workflows.WorkflowRunWithJobs[0]
	if len(selectedWorkflow.Runs) == 0 {
		return nil
	}

	currentIndex := m.Table.GetCurrItem()
	if currentIndex < 0 || currentIndex >= len(selectedWorkflow.Runs) {
		return nil
	}

	return selectedWorkflow.Runs[currentIndex]
}
