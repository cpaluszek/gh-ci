package commands

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
)

type ClientInitMsg struct {
	Client *github.Client
}

type ConfigInitMsg struct {
	Config *config.Config
}

type RepositoriesMsg struct {
	Repositories []*github.RepositoryData
}

type WorkflowsMsg struct {
	Workflows *github.RepositoryData
}

type SectionChangedMsg struct{}

type ErrorMsg struct {
	Error error
}

// Commands
// A command with no arguments
func InitConfig() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return ConfigInitMsg{
		Config: cfg,
	}
}

// A command with arguments
func InitClient(token string) tea.Cmd {
	return func() tea.Msg {
		client, err := github.NewClient(token)
		if err != nil {
			return ErrorMsg{
				Error: err,
			}
		}
		return ClientInitMsg{
			Client: client,
		}
	}
}

func SectionChanged() tea.Msg {
	return SectionChangedMsg{}
}

func FetchRepositories(client *github.Client, names []string) tea.Cmd {
	return func() tea.Msg {
		repos, err := client.FetchRepositoriesWithWorkflows(names)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return RepositoriesMsg{
			Repositories: repos,
		}
	}
}

func GoToWorkflow(row github.RowData) tea.Cmd {
	return func() tea.Msg {
		workflows, ok := row.(*github.RepositoryData)
		if !ok {
			return ErrorMsg{
				Error: fmt.Errorf("selected row is not a repository"),
			}
		}
		return WorkflowsMsg{
			Workflows: workflows,
		}
	}
}
