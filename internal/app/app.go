package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/google/go-github/v71/github"

	"github.com/cpaluszek/pipeye/internal/config"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/ui"
)

type Model struct {
	BaseView
	config       *config.Config
	Client       *github.Client
	repositories []*gh.Repository
	// Detail view
	selectedRepo  *gh.Repository
	selectedIndex int
	detailView    DetailView
	showDetail    bool
}

func New(cfg *config.Config) *Model {
	baseView := NewBaseView(viewport.New(0, 0), nil)
	return &Model{
		BaseView:      baseView,
		config:        cfg,
		selectedIndex: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		InitClient(m.config.Github.Token),
		m.Spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.UpdateSize(msg.Width, msg.Height)
		if m.showDetail {
			m.detailView.UpdateSize(msg.Width, msg.Height)
		}
		return m, nil

	case tea.KeyMsg:
		if m.showDetail {
			// Handle key events in detail view
			// TODO: esc feels slow because of terminal delay
			switch msg.String() {
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
				if !m.Loading && len(m.repositories) > 0 {
					m.selectedIndex = min(m.selectedIndex+1, len(m.repositories)-1)
					return m, nil
				}
			case "k", "up":
				if !m.Loading && len(m.repositories) > 0 {
					m.selectedIndex = max(m.selectedIndex-1, 0)
					return m, nil
				}
			case "enter":
				if len(m.repositories) > 0 {
					m.selectedRepo = m.repositories[m.selectedIndex]
					m.detailView = NewDetailView(m.selectedRepo, m.Viewport, m.Client)
					m.showDetail = true
					return m, m.detailView.Init()
				}
			}
		}

		var cmd tea.Cmd
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)

	case ClientInitializedMsg:
		m.Client = msg.Client
		cmds = append(cmds, FetchRepositories(m.Client))

	case RepositoriesMsg:
		m.Loading = false
		if msg.Error != nil {
			m.Error = msg.Error
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
		m.Loading = false
		m.Error = msg.Err
		return m, nil

	default:
		if m.Loading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.showDetail && m.detailView.Loading {
			// Forward spinner updates to detail view when it's visible and loading
			var cmd tea.Cmd
			m.detailView, cmd = m.detailView.Update(msg)
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

	if m.Error != nil {
		return fmt.Sprintf("Error: %v\n\n(press q to quit)", m.Error)
	}

	var content string

	if m.Loading {
		content = fmt.Sprintf("%s Fetching repositories...\n\n", m.Spinner.View())
	} else if len(m.repositories) > 0 {
		content = ui.RenderRepositoriesTable(m.repositories, m.selectedIndex, m.Viewport.Width)
	} else if m.Client != nil {
		content = "No repositories found.\n"
	}

	m.Viewport.SetContent(content)

	statusBar := ui.RenderStatusBar(m.Loading, len(m.repositories), m.StatusBarStyle)

	return fmt.Sprintf("%s\n%s", m.Viewport.View(), statusBar)
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
