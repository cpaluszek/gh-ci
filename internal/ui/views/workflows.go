package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/models"
	"github.com/cpaluszek/pipeye/internal/ui/render"
	gh "github.com/google/go-github/v71/github"
)

// TODO: fetch workflow with repos to eliminate loading time

type WorkflowsView struct {
	BaseView
	repository            *gh.Repository
	workflowsWithRuns     []*github.WorkflowWithRuns
	selectedWorkflowIndex int
	selectedRunIndex      int
}

func NewWorkflowsView(repo *gh.Repository, vp viewport.Model, client *github.Client) WorkflowsView {
	baseView := NewBaseView(vp, client)
	return WorkflowsView{
		BaseView:              baseView,
		repository:            repo,
		selectedWorkflowIndex: 0,
		selectedRunIndex:      0,
	}
}

func (d WorkflowsView) Init() tea.Cmd {
	d.Loading = true
	d.Error = nil

	// Fetch workflows with runs for the repository
	return tea.Batch(
		d.Spinner.Tick,
		models.FetchWorkflows(d.Client, d.repository.GetFullName()),
	)
}

func (d WorkflowsView) Update(msg tea.Msg) (WorkflowsView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case models.WorkflowsViewMsg:
		d.Loading = false
		if msg.Error != nil {
			d.Error = msg.Error
			return d, nil
		}
		d.workflowsWithRuns = msg.WorkflowsWithRuns
		return d, nil

	case tea.WindowSizeMsg:
		d.UpdateSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return d, nil // we'll handle this in the main model

		case "j", "down":
			if !d.Loading && len(d.workflowsWithRuns) > 0 {
				runs := d.workflowsWithRuns[d.selectedWorkflowIndex].Runs
				if len(runs) > 0 {
					d.selectedRunIndex = min(d.selectedRunIndex+1, len(runs)-1)
				}
				return d, nil
			}
		case "k", "up":
			if !d.Loading && len(d.workflowsWithRuns) > 0 {
				runs := d.workflowsWithRuns[d.selectedWorkflowIndex].Runs
				if len(runs) > 0 {
					d.selectedRunIndex = max(d.selectedRunIndex-1, 0)
				}
				return d, nil
			}
		}

		// Forward key messages to the viewport
		d.Viewport, cmd = d.Viewport.Update(msg)
		cmds = append(cmds, cmd)

	default:
		if d.Loading {
			var spinnerCmd tea.Cmd
			d.Spinner, spinnerCmd = d.Spinner.Update(msg)
			cmds = append(cmds, spinnerCmd)
		}
	}

	if len(cmds) > 0 {
		return d, tea.Batch(cmds...)
	}
	return d, cmd
}

func (d WorkflowsView) View() string {
	var content string

	if d.Loading {
		content = fmt.Sprintf("%s Loading workflows...\n\n", d.Spinner.View())
	} else {
		content = render.RenderWorkflowsView(d.repository, d.workflowsWithRuns, d.selectedRunIndex, d.Loading, d.Error)
	}

	d.Viewport.SetContent(content)
	statusBar := render.RenderWorkflowsStatusBar(d.Loading, d.StatusBarStyle)

	return fmt.Sprintf("%s\n%s", d.Viewport.View(), statusBar)
}
