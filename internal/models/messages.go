package models

import (
	"github.com/cpaluszek/pipeye/internal/github"
	gh "github.com/google/go-github/v71/github"
)

type (
	ErrMsg struct {
		Err error
	}

	ClientInitializedMsg struct {
		Client *github.Client
	}

	RepositoriesMsg struct {
		Repositories []*gh.Repository
		Error        error
	}

	DetailViewMsg struct {
		WorkflowsWithRuns []*github.WorkflowWithRuns
		Error             error
	}
)

func NewErrMsg(err error) ErrMsg {
	return ErrMsg{Err: err}
}

func NewClientInitializedMsg(client *github.Client) ClientInitializedMsg {
	return ClientInitializedMsg{Client: client}
}

func NewRepositoriesMsg(repos []*gh.Repository, err error) RepositoriesMsg {
	return RepositoriesMsg{
		Repositories: repos,
		Error:        err,
	}
}

func NewDetailViewMsg(workflows []*github.WorkflowWithRuns, err error) DetailViewMsg {
	return DetailViewMsg{
		WorkflowsWithRuns: workflows,
		Error:             err,
	}
}
