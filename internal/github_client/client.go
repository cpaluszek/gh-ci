package github_client

import (
	"context"
	"time"

	"github.com/google/go-github/v71/github"
)

type Client struct {
	Client *github.Client
}

func NewClient(token string) (*Client, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	return &Client{
		Client: client,
	}, nil
}

func (c *Client) FetchRepositories(ctx context.Context) ([]*github.Repository, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()
	}

	opt := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := c.Client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}
