package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/config"
	"github.com/cpaluszek/pipeye/internal/github_client"
	"github.com/google/go-github/v71/github"
)

type (
	errMsg struct {
		err error
	}

	clientInitializedMsg struct {
		client *github_client.Client
	}

	repositoriesMsg struct {
		Repositories []*github.Repository
		Error        error
	}
)

type Model struct {
	config *config.Config
	client *github_client.Client
	repositories []*github.Repository
	width int
	height int
	error error
	loading bool
	spinner spinner.Model
}

func New(cfg *config.Config) *Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Model{
		config: cfg,
		loading: true,
		spinner: s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.initClient,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case clientInitializedMsg:
		m.client = msg.client
		return m, m.fetchRepositories()

	case repositoriesMsg:
		m.loading = false
		if msg.Error != nil {
			m.error = msg.Error
			return m, nil
		}
		fmt.Printf("Received %d repositories in UI\n", len(msg.Repositories))
		m.repositories = msg.Repositories
		return m, nil
	
	case errMsg:
        m.loading = false
        m.error = msg.err
        return m, nil

	default:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	}

	return m, nil
}

func (m Model) View() string {
	if m.error != nil {
        return fmt.Sprintf("Error: %v\n\n(press q to quit)", m.error)
    }

	var sb strings.Builder

    if m.loading {
		sb.WriteString(fmt.Sprintf("%s Fetching repositories...\n\n", m.spinner.View()))
    } else {
		if len(m.repositories) > 0 {
			sb.WriteString(fmt.Sprintf("Found %d repositories with workflows:\n\n", len(m.repositories)))
			for i, repo := range m.repositories {
				sb.WriteString(fmt.Sprintf("%d. %s - %s (%d stars)\n", 
					i+1, 
					repo.GetFullName(), 
					repo.GetDescription(),
					repo.GetStargazersCount(),
					))
				if i >= 9 { // Show only 10 repos for now
					sb.WriteString(fmt.Sprintf("\n... and %d more\n", len(m.repositories)-10))
					break
				}
			}
		} else if m.client != nil {
			sb.WriteString("No repositories found.\n")
		}
	}

    sb.WriteString("\n(press q to quit)")
 
    return sb.String()
}

func (m Model) initClient() tea.Msg {
	client, err := github_client.NewClient(m.config.Github.Token)
	if err != nil {
		return errMsg{err};
	}

	return clientInitializedMsg{client: client}
}

func (m Model) fetchRepositories() tea.Cmd {
	return func() tea.Msg {
		repos, err := m.client.FetchRepositoriesWithWorkflows()
        if err != nil {
            return repositoriesMsg{Error: err}
        }
 
        return repositoriesMsg{
            Repositories: repos,
            Error: nil,
        }
	}
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
		tea.WithMouseCellMotion(),
		)

	_, error := p.Run()
	return error
}
