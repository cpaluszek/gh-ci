package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/ui/render"
)

type WorkflowsView struct {
	BaseView
	repository            *github.RepositoryData
	selectedWorkflowIndex int
	selectedRunIndex      int
}

func NewWorkflowsView(repo *github.RepositoryData, vp viewport.Model, client *github.Client) WorkflowsView {
	baseView := NewBaseView(vp, client, false)
	return WorkflowsView{
		BaseView:              baseView,
		repository:            repo,
		selectedWorkflowIndex: 0,
		selectedRunIndex:      0,
	}
}

func (d WorkflowsView) Init() tea.Cmd {
	d.Error = nil

	return nil
}

func (d WorkflowsView) Update(msg tea.Msg) (WorkflowsView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.UpdateSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return d, nil // we'll handle this in the main model

		case "j", "down":
			if len(d.repository.WorkflowRunWithJobs) > 0 {
				runs := d.repository.WorkflowRunWithJobs[d.selectedWorkflowIndex].Runs
				if len(runs) > 0 {
					d.selectedRunIndex = min(d.selectedRunIndex+1, len(runs)-1)
				}
				return d, nil
			}
		case "k", "up":
			if len(d.repository.WorkflowRunWithJobs) > 0 {
				runs := d.repository.WorkflowRunWithJobs[d.selectedWorkflowIndex].Runs
				if len(runs) > 0 {
					d.selectedRunIndex = max(d.selectedRunIndex-1, 0)
				}
				return d, nil
			}
		}

		// Forward key messages to the viewport
		d.Viewport, cmd = d.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	if len(cmds) > 0 {
		return d, tea.Batch(cmds...)
	}
	return d, cmd
}

func (d WorkflowsView) View() string {
	content := render.RenderWorkflowsView(d.repository, d.selectedRunIndex, d.Viewport.Width, d.Error)

	d.Viewport.SetContent(content)
	statusBar := render.RenderWorkflowsStatusBar(d.StatusBarStyle)

	return fmt.Sprintf("%s\n%s", d.Viewport.View(), statusBar)
}
