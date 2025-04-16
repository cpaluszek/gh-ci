package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gh "github.com/google/go-github/v71/github"

	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/models"
	"github.com/cpaluszek/pipeye/internal/ui/render"
)

type ListView struct {
	BaseView
	Repositories  []*gh.Repository
	selectedIndex int
}

func NewListView(client *github.Client) ListView {
	baseView := NewBaseView(viewport.New(0, 0), client)
	return ListView{
		BaseView:      baseView,
		selectedIndex: 0,
	}
}

func (l ListView) Init() tea.Cmd {
	l.Loading = true
	l.Error = nil

	// Fetch repositories
	return tea.Batch(
		models.FetchRepositories(l.Client),
		l.Spinner.Tick,
	)
}

func (l ListView) Update(msg tea.Msg) (ListView, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.UpdateSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if !l.Loading && len(l.Repositories) > 0 {
				l.selectedIndex = min(l.selectedIndex+1, len(l.Repositories)-1)
				return l, nil
			}
		case "k", "up":
			if !l.Loading && len(l.Repositories) > 0 {
				l.selectedIndex = max(l.selectedIndex-1, 0)
				return l, nil
			}
		}

		// Forward key messages to the viewport
		l.Viewport, cmd = l.Viewport.Update(msg)
		cmds = append(cmds, cmd)

	case models.RepositoriesMsg:
		l.Loading = false
		if msg.Error != nil {
			l.Error = msg.Error
			return l, nil
		}
		l.Repositories = msg.Repositories
		return l, nil

	default:
		if l.Loading {
			var spinnerCmd tea.Cmd
			l.Spinner, spinnerCmd = l.Spinner.Update(msg)
			cmds = append(cmds, spinnerCmd)
		}
	}

	if len(cmds) > 0 {
		return l, tea.Batch(cmds...)
	}
	return l, cmd
}

func (l ListView) View() string {
	if l.Error != nil {
		return fmt.Sprintf("Error: %v\n\n(press q to quit)", l.Error)
	}

	var content string

	if l.Loading {
		content = fmt.Sprintf("%s Fetching repositories...\n\n", l.Spinner.View())
	} else if len(l.Repositories) > 0 {
		content = render.RenderRepositoriesTable(l.Repositories, l.selectedIndex, l.Viewport.Width)
	} else if l.Client != nil {
		content = "No repositories found.\n"
	}

	l.Viewport.SetContent(content)
	statusBar := render.RenderStatusBar(l.Loading, len(l.Repositories), l.StatusBarStyle)

	return fmt.Sprintf("%s\n%s", l.Viewport.View(), statusBar)
}

func (l ListView) GetSelectedRepo() *gh.Repository {
	if l.selectedIndex >= 0 && l.selectedIndex < len(l.Repositories) {
		return l.Repositories[l.selectedIndex]
	}
	return nil
}

func (l ListView) HasRepositories() bool {
	return len(l.Repositories) > 0
}
