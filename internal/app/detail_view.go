package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/github_client"
	"github.com/cpaluszek/pipeye/internal/ui"
	"github.com/google/go-github/v71/github"
)

// TODO: fetch workflow with repos to eliminate loading time

type DetailViewMsg struct {
	Workflows []*github.Workflow
	Error     error
}

type DetailView struct {
	repository  *github.Repository
	workflows   []*github.Workflow
	loading     bool
	error       error
	viewport    viewport.Model
	client      *github_client.Client
	statusBarStyle lipgloss.Style
}

func NewDetailView(repo *github.Repository, client *github_client.Client, viewport viewport.Model) DetailView {
	return DetailView{
		repository: repo,
		loading:    true,
		client:     client,
		viewport:   viewport,
		statusBarStyle: ui.StatusStyle.Width(viewport.Width),
	}
}

func (d DetailView) FetchWorkflows() tea.Cmd {
    return func() tea.Msg {
        owner, repo := github_client.ParseFullName(*d.repository.FullName)
        workflows, err := d.client.FetchWorkflows(owner, repo)
        return DetailViewMsg{
            Workflows: workflows,
            Error:     err,
        }
    }
}

func (d DetailView) Update(msg tea.Msg) (DetailView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case DetailViewMsg:
		d.loading = false
		if msg.Error != nil {
			d.error = msg.Error
			return d, nil
		}
		d.workflows = msg.Workflows
		return d, nil

	case tea.WindowSizeMsg:
		d.viewport.Width = msg.Width
		d.viewport.Height = msg.Height - 1 // reserve space for status
		d.statusBarStyle = ui.StatusStyle.Width(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			return d, nil // we'll handle this in the main model
		}

		// Forward key messages to the viewport
		d.viewport, cmd = d.viewport.Update(msg)
		return d, cmd
	}

	return d, nil
}

func (d DetailView) View() string {
    content := ui.RenderDetailView(d.repository, d.workflows, d.loading, d.error)
    d.viewport.SetContent(content)
	statusBar := ui.RenderDetailViewStatusBar(d.loading, d.statusBarStyle)

	return fmt.Sprintf("%s\n%s", d.viewport.View(), statusBar)
}
