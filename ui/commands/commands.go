package commands

import (
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

func FetchRepositories(client *github.Client) tea.Cmd {
	return func() tea.Msg {
		repos, err := client.FetchRepositoriesWithWorkflows()
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return RepositoriesMsg{
			Repositories: repos,
		}
	}
}
