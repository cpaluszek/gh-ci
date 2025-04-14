package models

import (
	tea "github.com/charmbracelet/bubbletea"
	github "github.com/cpaluszek/pipeye/internal/github"
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

func FetchWorkflows(client *github.Client, repoName string) tea.Cmd {
	return func() tea.Msg {
		owner, repo := github.ParseFullName(repoName)
		workflowsWithRuns, err := client.FetchWorkflowsWithRuns(owner, repo)
		return NewDetailViewMsg(workflowsWithRuns, err)
	}
}
