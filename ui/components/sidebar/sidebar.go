package sidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/styles"
	"github.com/cpaluszek/pipeye/ui/utils"
)

type Model struct {
	ctx      *context.Context
	viewport viewport.Model
	data     string
}

func NewModel(ctx *context.Context) Model {
	return Model{
		data: "",
		viewport: viewport.Model{
			Width:  0,
			Height: 0,
		},
		ctx: ctx,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	height := m.ctx.MainContentHeight
	width := constants.SideBarWidth

	style := styles.SideBarStyle.
		Height(height).
		MaxHeight(height).
		Width(width).
		MaxWidth(width)

	if m.data == "" {
		return style.Align(lipgloss.Center).Render(
			lipgloss.PlaceVertical(height, lipgloss.Center, ""),
		)
	}

	return style.Render(m.viewport.View())
}

func (m *Model) SetContent(data string) {
	m.data = data
	m.viewport.SetContent(data)
}

func (m *Model) UpdateProgramContext(ctx *context.Context) {
	if ctx == nil {
		return
	}

	m.ctx = ctx
	m.viewport.Height = m.ctx.MainContentHeight
	m.viewport.Width = constants.SideBarWidth
}

func (m *Model) GenerateRepoSidebarContent(repo *github.Repository) {
	content := []string{
		styles.TitleStyle.Render("Repository: " + repo.GetName()),
		"",
	}

	// If no workflows, show message and return
	if len(repo.Workflows) == 0 || len(repo.Workflows[0].Runs) == 0 {
		content = append(content, styles.DefaultStyle.Render("No workflows found"))
		m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
		return
	}

	workflowDisplayHeight := 5
	for i, workflow := range repo.Workflows {
		if len(content) >= m.viewport.Height-workflowDisplayHeight {
			content = append(content, styles.DefaultStyle.Render(fmt.Sprintf("\n+ %d more workflows...", len(repo.Workflows)-i)))
			break
		}

		if len(workflow.Runs) == 0 {
			continue
		}

		latestRun := workflow.Runs[0].Info

		workflowName := utils.TruncateString(*workflow.Info.Name, constants.SideBarWidth-4)

		createdTime := styles.DefaultStyle.Render(utils.FormatTime(latestRun.GetCreatedAt().Time))
		statusDuration := utils.GetWorkflowRunStatus(latestRun) + " · " + createdTime

		commitMsg := strings.Split(latestRun.GetHeadCommit().GetMessage(), "\n")[0]

		eventIcon := utils.GetRunEventSymbol(*latestRun.Event)

		content = append(content, styles.TitleStyle.Render(workflowName))
		content = append(content, styles.DefaultStyle.Render(statusDuration))
		content = append(content, styles.DefaultStyle.Render(eventIcon+latestRun.GetEvent()+" · "+commitMsg))

		if i < len(repo.Workflows)-1 {
			content = append(content, styles.DefaultStyle.Render(""))
		}
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (m *Model) GenerateWorkflowSidebarContent(workflow *github.WorkflowRun) {
	content := []string{
		styles.TitleStyle.Render("Workflow: " + workflow.GetName()),
		"",
	}

	if len(workflow.Jobs) > 0 {
		for _, job := range workflow.Jobs {
			content = append(content, styles.TitleStyle.Render(*job.Name))
			status := utils.GetJobStatusSymbol(job.GetStatus(), job.GetConclusion())
			status = utils.CleanANSIEscapes(status) + job.GetConclusion()
			content = append(content, styles.DefaultStyle.Render(status))
			if job.GetConclusion() != "skipped" {
				content = append(content, styles.DefaultStyle.Render("Duration: "+utils.GetJobDuration(job)))
			}
			content = append(content, "")
		}
	} else {
		content = append(content, styles.DefaultStyle.Render("No jobs found"))
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (m *Model) GenerateRunSidebarContent(run *github.Job) {
	content := []string{
		styles.TitleStyle.Render("Run: " + run.GetName()),
		"",
	}

	for _, step := range run.Job.Steps {
		if step.GetName() == "" {
			continue
		}
		content = append(content, styles.TitleStyle.Render(*step.Name))
		content = append(content, "")
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}
