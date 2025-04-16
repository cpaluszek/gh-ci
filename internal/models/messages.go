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

	WorkflowsViewMsg struct {
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

func NewWorkflowsViewMsg(workflows []*github.WorkflowWithRuns, err error) WorkflowsViewMsg {
	return WorkflowsViewMsg{
		WorkflowsWithRuns: workflows,
		Error:             err,
	}
}
