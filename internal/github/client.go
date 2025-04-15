package github

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	gh "github.com/google/go-github/v71/github"
)

type Client struct {
	Client *gh.Client
}

type WorkflowWithRuns struct {
	Workflow *gh.Workflow
	Runs     []*WorkflowRunWithJobs
	Error    error
}

type WorkflowRunWithJobs struct {
	Run  *gh.WorkflowRun
	Jobs []*gh.WorkflowJob
}

const (
	defaultConcurrency = 10
	defaultTimeout     = 10 * time.Second
)

// TODO: if fetching is used on interval, should use cache for repo workflows

func NewClient(token string) (*Client, error) {
	client := gh.NewClient(nil).WithAuthToken(token)
	return &Client{
		Client: client,
	}, nil
}

// FetchRepositoriesWithWorkflows fetches repositories that have GitHub Actions workflows
func (c *Client) FetchRepositoriesWithWorkflows() ([]*gh.Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	repos, err := c.fetchRepositories(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	log.Printf("Fetching workflows for %d repositories...\n", len(repos))

	var result []*gh.Repository
	var wg sync.WaitGroup
	var mutex sync.Mutex
	semaphore := make(chan struct{}, defaultConcurrency)

	for _, repo := range repos {
		wg.Add(1)
		currentRepo := repo

		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			workflows, _, err := c.Client.Actions.ListWorkflows(
				ctx,
				currentRepo.GetOwner().GetLogin(),
				currentRepo.GetName(),
				&gh.ListOptions{PerPage: 1},
			)

			if err != nil {
				log.Printf("Error fetching workflows for %s: %v",
					currentRepo.GetFullName(), err)
				return
			}

			if workflows != nil && workflows.GetTotalCount() > 0 {
				mutex.Lock()
				result = append(result, currentRepo)
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	log.Printf("Found %d repositories with workflows in %s", len(result), time.Since(start))
	return result, nil
}

// FetchWorkflowsWithRuns fetches workflows and their recent runs for a repository
func (c *Client) FetchWorkflowsWithRuns(owner, repo string) ([]*WorkflowWithRuns, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	log.Printf("Fetching workflows and runs for %s/%s...\n", owner, repo)

	// Fetch workflows
	workflows, _, err := c.Client.Actions.ListWorkflows(ctx, owner, repo, &gh.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Found %d workflows", len(workflows.Workflows))

	var result []*WorkflowWithRuns
	var wgWorkflows sync.WaitGroup
	var mutexWorkflows sync.Mutex
	semaphoreWorkflows := make(chan struct{}, defaultConcurrency)

	// Process workflows in parallel
	for _, workflow := range workflows.Workflows {
		wgWorkflows.Add(1)
		currentWorkflow := workflow

		go func() {
			defer wgWorkflows.Done()

			// Acquire semaphore slot
			semaphoreWorkflows <- struct{}{}
			defer func() { <-semaphoreWorkflows }()

			workflowWithRuns := &WorkflowWithRuns{
				Workflow: currentWorkflow,
			}

			// Fetch runs for this workflow
			runs, _, err := c.Client.Actions.ListWorkflowRunsByID(
				ctx,
				owner,
				repo,
				currentWorkflow.GetID(),
				&gh.ListWorkflowRunsOptions{
					ListOptions: gh.ListOptions{PerPage: 20},
				},
			)

			if err != nil {
				workflowWithRuns.Error = err
				mutexWorkflows.Lock()
				result = append(result, workflowWithRuns)
				mutexWorkflows.Unlock()
				return
			}

			// Fetch jobs for each run
			if runs != nil && len(runs.WorkflowRuns) > 0 {
				runsWithJobs := make([]*WorkflowRunWithJobs, len(runs.WorkflowRuns))
				var wgJobs sync.WaitGroup
				var mutexJobs sync.Mutex

				for i, run := range runs.WorkflowRuns {
					wgJobs.Add(1)
					runIndex := i
					currentRun := run

					go func() {
						defer wgJobs.Done()

						runWithJobs := &WorkflowRunWithJobs{
							Run: currentRun,
						}

						jobs, _, err := c.Client.Actions.ListWorkflowJobs(
							ctx,
							owner,
							repo,
							currentRun.GetID(),
							&gh.ListWorkflowJobsOptions{
								ListOptions: gh.ListOptions{PerPage: 100},
							},
						)

						if err == nil && jobs != nil {
							runWithJobs.Jobs = jobs.Jobs
						}

						mutexJobs.Lock()
						runsWithJobs[runIndex] = runWithJobs
						mutexJobs.Unlock()
					}()
				}

				wgJobs.Wait()
				workflowWithRuns.Runs = runsWithJobs
			}

			mutexWorkflows.Lock()
			result = append(result, workflowWithRuns)
			mutexWorkflows.Unlock()
		}()
	}

	wgWorkflows.Wait()
	log.Printf("Fetched %d workflows with runs in %s", len(result), time.Since(start))
	return result, nil
}

// Helper function to fetch repositories
func (c *Client) fetchRepositories(ctx context.Context) ([]*gh.Repository, error) {
	opt := &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{PerPage: 20},
	}

	var allRepos []*gh.Repository
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

// ParseFullName splits a full repository name into owner and repo parts
func ParseFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
