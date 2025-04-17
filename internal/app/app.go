package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cpaluszek/pipeye/internal/config"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/models"
	"github.com/cpaluszek/pipeye/internal/ui/views"
)

type Model struct {
	views.BaseView
	config           *config.Config
	Client           *github.Client
	RepositoriesView views.RepositoriesView
	workflowsView    views.WorkflowsView
	showWorkflows    bool
}

func New(cfg *config.Config) *Model {
	return &Model{
		config:        cfg,
		showWorkflows: false,
	}
}

func (m Model) Init() tea.Cmd {
	return models.InitClient(m.config.Github.Token)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.showWorkflows {
			workflowsView, cmd := m.workflowsView.Update(msg)
			m.workflowsView = workflowsView
			cmds = append(cmds, cmd)
		} else {
			repositoryView, cmd := m.RepositoriesView.Update(msg)
			m.RepositoriesView = repositoryView
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if m.showWorkflows {
			// Handle key events in workflows view
			switch msg.String() {
			case "esc", "backspace":
				m.showWorkflows = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.workflowsView, cmd = m.workflowsView.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				if m.RepositoriesView.HasRepositories() {
					selectedRepo := m.RepositoriesView.GetSelectedRepo()
					m.workflowsView = views.NewWorkflowsView(selectedRepo, m.RepositoriesView.Viewport, m.Client)
					m.workflowsView.UpdateSize(m.RepositoriesView.Viewport.Width, m.RepositoriesView.Viewport.Height)
					m.showWorkflows = true
					return m, m.workflowsView.Init()
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.RepositoriesView, cmd = m.RepositoriesView.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case models.ClientInitializedMsg:
		m.Client = msg.Client
		m.RepositoriesView = views.NewRepositoriesView(m.Client)
		return m, m.RepositoriesView.Init()

	case models.RepositoriesMsg:
		var cmd tea.Cmd
		m.RepositoriesView, cmd = m.RepositoriesView.Update(msg)
		cmds = append(cmds, cmd)

	case models.WorkflowsViewMsg:
		if m.showWorkflows {
			var cmd tea.Cmd
			m.workflowsView, cmd = m.workflowsView.Update(msg)
			cmds = append(cmds, cmd)
		}

	case models.ErrMsg:
		if m.showWorkflows {
			m.workflowsView.Error = msg.Err
			m.workflowsView.Loading = false
		} else {
			m.RepositoriesView.Error = msg.Err
			m.RepositoriesView.Loading = false
		}
		return m, nil

	default:
		if !m.showWorkflows {
			var cmd tea.Cmd
			m.RepositoriesView, cmd = m.RepositoriesView.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.workflowsView, cmd = m.workflowsView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() string {
	if m.showWorkflows {
		return m.workflowsView.View()
	}
	return m.RepositoriesView.View()
}

func (m *Model) Run() error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	_, error := p.Run()
	return error
}
