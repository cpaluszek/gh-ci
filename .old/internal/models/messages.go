package models

import (
	"github.com/cpaluszek/pipeye/internal/github"
)

type (
	ErrMsg struct {
		Err error
	}

	ClientInitializedMsg struct {
		Client *github.Client
	}

	RepositoriesMsg struct {
		Repositories []*github.RepositoryData
		Error        error
	}

	WorkflowsViewMsg struct {
		Repository *github.RepositoryData
		Error      error
	}
)

func NewErrMsg(err error) ErrMsg {
	return ErrMsg{Err: err}
}

func NewClientInitializedMsg(client *github.Client) ClientInitializedMsg {
	return ClientInitializedMsg{Client: client}
}

func NewRepositoriesMsg(repos []*github.RepositoryData, err error) RepositoriesMsg {
	return RepositoriesMsg{
		Repositories: repos,
		Error:        err,
	}
}
