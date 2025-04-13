package github_client

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v71/github"
)

type Client struct {
	Client *github.Client
}

// TODO: if fetching is used on interval, should use cache for repo workflows

func NewClient(token string) (*Client, error) {
	client := github.NewClient(nil).WithAuthToken(token)
	return &Client{
		Client: client,
	}, nil
}

func (c *Client) fetchRepositories() ([]*github.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

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

func (c *Client) FetchRepositoriesWithWorkflows() ([]*github.Repository, error) {
	repos, err := c.fetchRepositories()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	var repositoryWorkflows []*github.Repository
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Process repositories in parallel with a limit
	semaphore := make(chan struct{}, 5)

	for _, repo := range repos {
		wg.Add(1)

		go func(repo *github.Repository) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <- semaphore }()

			workflows, _, err := c.Client.Actions.ListWorkflows(
				ctx,
				repo.GetOwner().GetLogin(),
				repo.GetName(),
				&github.ListOptions{PerPage: 1},
			)

			if err != nil {
				return
			}

			if workflows.GetTotalCount() > 0 {
				mutex.Lock()
				repositoryWorkflows = append(repositoryWorkflows, repo)
				mutex.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	return repositoryWorkflows, nil
}

func (c *Client) FetchWorkflows(owner, repo string) ([]*github.Workflow, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// TODO: define workflow count
	worflows, _, err := c.Client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{
		PerPage: 20,
	})
	if err != nil {
		return nil, err
	}

	return worflows.Workflows, nil
}

// ParseFullName splits a full repository name into owner and repo parts
func ParseFullName(fullName string) (string, string) {
    parts := strings.Split(fullName, "/")
    if len(parts) == 2 {
        return parts[0], parts[1]
    }
    return "", ""
}
