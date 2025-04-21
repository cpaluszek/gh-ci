package commands

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
)

type ClientInitMsg struct {
	Client *github.Client
	Error  error
}

type ConfigInitMsg struct {
	Config *config.Config
	Error  error
}

// Commands
// A command with no arguments
func InitConfig() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return ConfigInitMsg{
			Config: nil,
			Error:  err,
		}
	}
	return ConfigInitMsg{
		Config: cfg,
		Error:  nil,
	}
}

// A command with arguments
func InitClient(token string) tea.Cmd {
	return func() tea.Msg {
		client, err := github.NewClient(token)
		if err != nil {
			return ClientInitMsg{
				Client: nil,
				Error:  err,
			}
		}
		return ClientInitMsg{
			Client: client,
			Error:  nil,
		}
	}
}
