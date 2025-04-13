package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/viewport"
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

// TODO: move to styles
var (
    statusStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(lipgloss.Color("#333333")).
        Padding(0, 1).
        Bold(false)
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
		statusBarStyle: statusStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.initClient,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 1 // reserve space for status
		m.statusBarStyle = statusStyle.Width(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		// case "up", "k":
		// 	m.viewport.ScrollUp(1)
		// case "down", "j":
		// 	m.viewport.ScrollDown(1)
		// case "pgup":
		// 	m.viewport.HalfPageUp()
		// case "pgdown":
		// 	m.viewport.HalfPageDown()
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case clientInitializedMsg:
		m.client = msg.client
		cmds = append(cmds, m.fetchRepositories())

	case repositoriesMsg:
		m.loading = false
		if msg.Error != nil {
			m.error = msg.Error
			return m, nil
		}
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
			cmds = append(cmds, cmd)
		}

	}

	if len(cmds) > 0 {
        return m, tea.Batch(cmds...)
    }

	return m, nil
}

func (m Model) View() string {
	if m.error != nil {
        return fmt.Sprintf("Error: %v\n\n(press q to quit)", m.error)
    }

	var s strings.Builder

	 // Calculate dynamic widths based on available space
    availableWidth := m.viewport.Width
    if availableWidth == 0 {
        availableWidth = 100 // Fallback if window size not detected yet
    }

    nameWidth := int(float64(availableWidth) * 0.4)
    langWidth := int(float64(availableWidth) * 0.15)
    starsWidth := int(float64(availableWidth) * 0.10)
    updatedWidth := int(float64(availableWidth) * 0.25)
    workflowsWidth := int(float64(availableWidth) * 0.10)

    // Ensure minimum widths
    nameWidth = max(nameWidth, 20)
    langWidth = max(langWidth, 10)
    starsWidth = max(starsWidth, 6)
    updatedWidth = max(updatedWidth, 15)
    workflowsWidth = max(workflowsWidth, 8)

    // Adjust if total exceeds available width
    totalWidth := nameWidth + langWidth + starsWidth + updatedWidth + workflowsWidth
    if totalWidth > availableWidth {
        // Reduce proportionally
        ratio := float64(availableWidth) / float64(totalWidth)
        nameWidth = int(float64(nameWidth) * ratio)
        langWidth = int(float64(langWidth) * ratio)
        starsWidth = int(float64(starsWidth) * ratio)
        updatedWidth = int(float64(updatedWidth) * ratio)
        workflowsWidth = int(float64(workflowsWidth) * ratio)
    }

    totalWidth = nameWidth + langWidth + starsWidth + updatedWidth + workflowsWidth

	 if m.loading {
        s.WriteString(fmt.Sprintf("%s Fetching repositories...\n\n", m.spinner.View()))
    } else {
        if len(m.repositories) > 0 {
            s.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#76ABDF")).Render("GitHub Repositories with Workflows"))
            s.WriteString("\n\n")

            // Column headers
            headers := lipgloss.JoinHorizontal(lipgloss.Top,
                lipgloss.NewStyle().Width(nameWidth).Bold(true).Align(lipgloss.Left).Render("Repository"),
                lipgloss.NewStyle().Width(langWidth).Bold(true).Align(lipgloss.Left).Render("Language"),
                lipgloss.NewStyle().Width(starsWidth).Bold(true).Align(lipgloss.Left).Render("Stars"),
                lipgloss.NewStyle().Width(updatedWidth).Bold(true).Align(lipgloss.Left).Render("Last Updated"),
                lipgloss.NewStyle().Width(workflowsWidth).Bold(true).Align(lipgloss.Left).Render("Workflows"),
            )
            s.WriteString(headers + "\n")
            s.WriteString(strings.Repeat("─", totalWidth) + "\n")

            for i, repo := range m.repositories {
                language := ""
                if repo.Language != nil {
                    language = *repo.Language
                }

                stars := "0"
                if repo.StargazersCount != nil {
                    stars = fmt.Sprintf("%d", *repo.StargazersCount)
                }

                updated := "Unknown"
                if repo.UpdatedAt != nil {
                    updated = repo.UpdatedAt.Format("Jan 2, 2006")
                }

                // Create row style with alternating background for better readability
                rowStyle := lipgloss.NewStyle()
                if i % 2 == 1 {
                    rowStyle = rowStyle.Background(lipgloss.Color("0"))
                }

                row := lipgloss.JoinHorizontal(lipgloss.Top,
                    rowStyle.Width(nameWidth).Align(lipgloss.Left).Render(*repo.FullName),
                    rowStyle.Width(langWidth).Align(lipgloss.Left).Render(language),
                    rowStyle.Width(starsWidth).Align(lipgloss.Left).Render(stars),
                    rowStyle.Width(updatedWidth).Align(lipgloss.Left).Render(updated),
                    rowStyle.Width(workflowsWidth).Align(lipgloss.Left).Render("✓"),
                )
                s.WriteString(row + "\n")
            }
        } else if m.client != nil {
            s.WriteString("No repositories found.\n")
        }
    }

    m.viewport.SetContent(s.String())	

    statusContent := ""
    if m.loading {
        statusContent = "Loading repositories... "
    } else if len(m.repositories) > 0 {
        statusContent = fmt.Sprintf("Found %d repositories with workflows · q: quit · ↑/↓: scroll", len(m.repositories))
    } else {
        statusContent = "No repositories found · q: quit"
    }

    return fmt.Sprintf("%s\n%s", m.viewport.View(), m.statusBarStyle.Render(statusContent))
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
		// tea.WithMouseCellMotion(),
		)

	_, error := p.Run()
	return error
}
