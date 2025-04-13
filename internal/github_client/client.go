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

type WorkflowWithRuns struct {
	Workflow *github.Workflow
	Runs     []*github.WorkflowRun
	Error    error
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

func (c *Client) FetchWorkflowsWithRuns(owner, repo string) ([]*WorkflowWithRuns, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Fetch workflows
	workflows, _, err := c.Client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}

	var result []*WorkflowWithRuns

	// For each workflow, fetch recent runs
	for _, workflow := range workflows.Workflows {
		workflowWithRuns := &WorkflowWithRuns{
			Workflow: workflow,
		}

		// Fetch 5 most recent runs for this workflow
		runs, _, err := c.Client.Actions.ListWorkflowRunsByID(
			ctx, 
			owner, 
			repo, 
			workflow.GetID(), 
			&github.ListWorkflowRunsOptions{
				ListOptions: github.ListOptions{PerPage: 20},
			},
			)

		if err != nil {
			workflowWithRuns.Error = err
		} else if runs != nil {
			workflowWithRuns.Runs = runs.WorkflowRuns
		}

		result = append(result, workflowWithRuns)
	}

	return result, nil
}

// ParseFullName splits a full repository name into owner and repo parts
func ParseFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
