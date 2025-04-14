package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/ui"
	gh "github.com/google/go-github/v71/github"
)

// TODO: fetch workflow with repos to eliminate loading time

type DetailView struct {
	BaseView
	repository        *gh.Repository
	workflowsWithRuns []*github.WorkflowWithRuns
}

func NewDetailView(repo *gh.Repository, vp viewport.Model, client *github.Client) DetailView {
	baseView := NewBaseView(vp, client)
	return DetailView{
		BaseView:   baseView,
		repository: repo,
	}
}

func (d DetailView) Init() tea.Cmd {
	d.Loading = true
	d.Error = nil

	// Fetch workflows with runs for the repository
	return tea.Batch(
		d.Spinner.Tick,
		FetchWorkflows(d.Client, d.repository.GetFullName()),
	)
}

func (d DetailView) Update(msg tea.Msg) (DetailView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case DetailViewMsg:
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

func (d DetailView) View() string {
	var content string

	if d.Loading {
		content = fmt.Sprintf("%s Loading workflows...\n\n", d.Spinner.View())
	} else {
		content = ui.RenderDetailView(d.repository, d.workflowsWithRuns, d.Loading, d.Error)
	}

	d.Viewport.SetContent(content)
	statusBar := ui.RenderDetailViewStatusBar(d.Loading, d.StatusBarStyle)

	return fmt.Sprintf("%s\n%s", d.Viewport.View(), statusBar)
}
