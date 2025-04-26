package github

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	gh "github.com/google/go-github/v71/github"
)

type Client struct {
	Client *gh.Client
}

const (
	defaultConcurrency = 10
	defaultTimeout     = 10 * time.Second
	workflowsPerPage   = 20
	workflowsRunCount  = 20
	jobsPerPage        = 10
)

func NewClient(token string) (*Client, error) {
	client := gh.NewClient(nil).WithAuthToken(token)
	return &Client{
		Client: client,
	}, nil
}

// FetchRepositoriesWithWorkflows fetches repositories that have GitHub Actions workflows
func (c *Client) FetchRepositoriesWithWorkflows(names []string) ([]*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// TODO: proper error management
	repos, err := c.fetchRepositories(ctx, names)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	log.Printf("Fetching workflows for %d repositories...\n", len(repos))

	var result []*Repository
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

			owner, repo := ParseFullName(currentRepo.GetFullName())
			workflowsWithRuns, error := c.FetchWorkflowsWithRuns(owner, repo)

			if error != nil {
				log.Printf("Error fetching workflows with runs for %s: %v",
					currentRepo.GetFullName(), error)
				return
			}
			mutex.Lock()
			result = append(result, &Repository{currentRepo, workflowsWithRuns, nil})
			mutex.Unlock()
		}()
	}

	wg.Wait()
	log.Printf("Found %d repositories with workflows in %s", len(result), time.Since(start))
	sort.Slice(result, func(i, j int) bool {
		return result[i].Info.UpdatedAt.After(result[j].Info.UpdatedAt.Time)
	})
	return result, nil
}

// FetchWorkflowsWithRuns fetches workflows and their recent runs for a repository
func (c *Client) FetchWorkflowsWithRuns(owner, repo string) ([]*Workflow, error) {
	// NOTE: use a new context?
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch workflows
	workflows, _, err := c.Client.Actions.ListWorkflows(ctx, owner, repo, &gh.ListOptions{
		PerPage: workflowsPerPage,
	})
	if err != nil {
		return nil, err
	}

	var result []*Workflow
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

			workflowWithRuns := &Workflow{
				Info: currentWorkflow,
			}

			// Fetch runs for this workflow
			runs, _, err := c.Client.Actions.ListWorkflowRunsByID(
				ctx,
				owner,
				repo,
				currentWorkflow.GetID(),
				&gh.ListWorkflowRunsOptions{
					ListOptions: gh.ListOptions{PerPage: workflowsRunCount},
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
				runsWithJobs := make([]*WorkflowRun, len(runs.WorkflowRuns))
				var wgJobs sync.WaitGroup
				var mutexJobs sync.Mutex

				for i, run := range runs.WorkflowRuns {
					wgJobs.Add(1)
					runIndex := i
					currentRun := run

					go func() {
						defer wgJobs.Done()

						runWithJobs := &WorkflowRun{
							Info: currentRun,
						}

						jobs, _, err := c.Client.Actions.ListWorkflowJobs(
							ctx,
							owner,
							repo,
							currentRun.GetID(),
							&gh.ListWorkflowJobsOptions{
								ListOptions: gh.ListOptions{PerPage: jobsPerPage},
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
	return result, nil
}

// Helper function to fetch repositories
func (c *Client) fetchRepositories(ctx context.Context, names []string) ([]*gh.Repository, error) {
	var repos []*gh.Repository

	// Fetch repositories from config
	for _, repoName := range names {
		repoParts := strings.Split(repoName, "/")
		if len(repoParts) != 2 {
			// TODO: format should be checked in config
			log.Printf("Invalid repository format: %s (expected 'owner/repo')", repoName)
			continue
		}
		repo, _, err := c.Client.Repositories.Get(ctx, repoParts[0], repoParts[1])
		if err != nil {
			log.Printf("Error fetching %s: %s", repoName, err)
			continue
		}

		repos = append(repos, repo)
	}

	return repos, nil
}

// ParseFullName splits a full repository name into owner and repo parts
func ParseFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
