package app

import (
    "github.com/cpaluszek/pipeye/internal/github_client"
    "github.com/google/go-github/v71/github"
)

type (
    ErrMsg struct {
        Err error
    }

    ClientInitializedMsg struct {
        Client *github_client.Client
    }

    RepositoriesMsg struct {
        Repositories []*github.Repository
        Error        error
    }
)

func NewErrMsg(err error) ErrMsg {
    return ErrMsg{Err: err}
}

func NewClientInitializedMsg(client *github_client.Client) ClientInitializedMsg {
    return ClientInitializedMsg{Client: client}
}

func NewRepositoriesMsg(repos []*github.Repository, err error) RepositoriesMsg {
    return RepositoriesMsg{
        Repositories: repos,
        Error:        err,
    }
}

