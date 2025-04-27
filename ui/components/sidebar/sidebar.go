package sidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/utils"
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

	style := m.ctx.Styles.SideBar.
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
		m.ctx.Styles.Title.Render("Repository: " + repo.GetName()),
		"",
	}

	// If no workflows, show message and return
	if len(repo.Workflows) == 0 || len(repo.Workflows[0].Runs) == 0 {
		content = append(content, m.ctx.Styles.Default.Render("No workflows found"))
		m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
		return
	}

	workflowDisplayHeight := 5
	for i, workflow := range repo.Workflows {
		if len(content) >= m.viewport.Height-workflowDisplayHeight {
			content = append(content, m.ctx.Styles.Default.Render(fmt.Sprintf("\n+ %d more workflows...", len(repo.Workflows)-i)))
			break
		}

		if len(workflow.Runs) == 0 {
			continue
		}

		latestRun := workflow.Runs[0]

		workflowName := utils.TruncateString(workflow.Name, constants.SideBarWidth-4)

		createdTime := m.ctx.Styles.Default.Render(utils.FormatTime(latestRun.CreatedAt))
		statusDuration := utils.GetWorkflowRunStatus(m.ctx, latestRun) + " · " + createdTime

		commitMsg := strings.Split(latestRun.HeadCommit.Message, "\n")[0]

		eventIcon := utils.GetRunEventSymbol(m.ctx, latestRun.Event)

		content = append(content, m.ctx.Styles.Title.Render(workflowName))
		content = append(content, m.ctx.Styles.Default.Render(statusDuration))
		content = append(content, m.ctx.Styles.Default.Render(eventIcon+latestRun.Event+" · "+commitMsg))

		if i < len(repo.Workflows)-1 {
			content = append(content, m.ctx.Styles.Default.Render(""))
		}
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (m *Model) GenerateWorkflowSidebarContent(workflow *github.WorkflowRun) {
	content := []string{
		m.ctx.Styles.Title.Render("Workflow: " + workflow.GetName()),
		"",
	}

	if len(workflow.Jobs) > 0 {
		for _, job := range workflow.Jobs {
			content = append(content, m.ctx.Styles.Title.Render(job.Name))
			status := utils.GetJobStatusSymbol(m.ctx, job.Status, job.Conclusion)
			status = utils.CleanANSIEscapes(status) + job.Conclusion
			content = append(content, m.ctx.Styles.Default.Render(status))
			if job.Conclusion != "skipped" {
				content = append(content, m.ctx.Styles.Default.Render("Duration: "+utils.GetJobDuration(job)))
			}
			content = append(content, "")
		}
	} else {
		content = append(content, m.ctx.Styles.Default.Render("No jobs found"))
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}

func (m *Model) GenerateRunSidebarContent(run *github.Job) {
	content := []string{
		m.ctx.Styles.Title.Render("Run: " + run.GetName()),
		"",
	}

	for _, step := range run.Steps {
		if step.Name == "" {
			continue
		}
		content = append(content, m.ctx.Styles.Title.Render(step.Name))
		status := utils.GetJobStatusSymbol(m.ctx, step.Status, step.Conclusion) + "· " + step.Conclusion
		content = append(content, m.ctx.Styles.Default.Render("Status: "+status))
		content = append(content, "")
	}

	m.SetContent(lipgloss.JoinVertical(lipgloss.Left, content...))
}
