package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v71/github"

	"github.com/cpaluszek/pipeye/internal/config"
	"github.com/cpaluszek/pipeye/internal/github_client"
	"github.com/cpaluszek/pipeye/internal/ui"
)

type Model struct {
	config *config.Config
	client *github_client.Client
	repositories []*github.Repository
	viewport viewport.Model
	error error
	loading bool
	spinner spinner.Model
	statusBarStyle lipgloss.Style
	// Detail view
	selectedRepo *github.Repository
	selectedIndex int
	detailView DetailView
	showDetail bool
}

func New(cfg *config.Config) *Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Model{
		config: cfg,
		loading: true,
		spinner: s,
		viewport: viewport.New(0, 0),
		statusBarStyle: ui.StatusStyle,
		selectedIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		InitClient(m.config.Github.Token),
		m.spinner.Tick,
		)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1 // reserve space for status
		m.statusBarStyle = ui.StatusStyle.Width(msg.Width)
		return m, nil

	case tea.KeyMsg:
		 if m.showDetail {
            // Handle key events in detail view
            switch msg.String() {
			// TODO: esc feels slow because of terminal delay
            case "esc", "backspace":
                m.showDetail = false
                return m, nil
            default:
                var cmd tea.Cmd
                m.detailView, cmd = m.detailView.Update(msg)
                cmds = append(cmds, cmd)
            }
        } else {
            // Handle key events in list view
            switch msg.String() {
            case "ctrl+c", "q":
                return m, tea.Quit
			case "j", "down":
				if !m.loading && len(m.repositories) > 0 {
					m.selectedIndex = min(m.selectedIndex+1, len(m.repositories)-1)
					return m, nil
				}
			case "k", "up":
				if !m.loading && len(m.repositories) > 0 {
					m.selectedIndex = max(m.selectedIndex-1, 0)
					return m, nil
				}
			case "enter":
				if len(m.repositories) > 0 {
					m.selectedRepo = m.repositories[m.selectedIndex]
					m.detailView = NewDetailView(m.selectedRepo, m.client, m.viewport)
					m.showDetail = true
					return m, m.detailView.FetchWorkflows()
				}
            }
        }

		var cmd tea.Cmd
        m.viewport, cmd = m.viewport.Update(msg)
        cmds = append(cmds, cmd)

	case ClientInitializedMsg:
		m.client = msg.Client
		cmds = append(cmds, FetchRepositories(m.client))

	case RepositoriesMsg:
		m.loading = false
		if msg.Error != nil {
			m.error = msg.Error
			return m, nil
		}
		m.repositories = msg.Repositories
		return m, nil

	case DetailViewMsg:
        if m.showDetail {
            var cmd tea.Cmd
            m.detailView, cmd = m.detailView.Update(msg)
            cmds = append(cmds, cmd)
        }

	case ErrMsg:
		m.loading = false
		m.error = msg.Err
		return m, nil

	default:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() string {
	if m.showDetail && m.selectedRepo != nil {
        return m.detailView.View()
    }

	if m.error != nil {
		return fmt.Sprintf("Error: %v\n\n(press q to quit)", m.error)
	}

	var content string

	if m.loading {
		content = fmt.Sprintf("%s Fetching repositories...\n\n", m.spinner.View())
	} else if len(m.repositories) > 0 {
		content = ui.RenderRepositoriesTable(m.repositories, m.selectedIndex, m.viewport.Width)
	} else if m.client != nil {
		content = "No repositories found.\n"
	}

	m.viewport.SetContent(content)

	statusBar := ui.RenderStatusBar(m.loading, len(m.repositories), m.statusBarStyle)

	return fmt.Sprintf("%s\n%s", m.viewport.View(), statusBar)
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
		*m,
		tea.WithAltScreen(),
		)

	_, error := p.Run()
	return error
}
