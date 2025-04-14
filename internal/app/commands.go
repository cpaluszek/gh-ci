package app

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/cpaluszek/pipeye/internal/github_client"
)

func InitClient(token string) tea.Cmd {
    return func() tea.Msg {
        client, err := github.NewClient(token)
        if err != nil {
            return ErrMsg{Err: err}
        }
        return ClientInitializedMsg{Client: client}
    }
}

func FetchRepositories(client *github.Client) tea.Cmd {
    return func() tea.Msg {
        repos, err := client.FetchRepositoriesWithWorkflows()
        if err != nil {
            return RepositoriesMsg{Error: err}
        }
        return RepositoriesMsg{
            Repositories: repos,
            Error:        nil,
        }
    }
}

