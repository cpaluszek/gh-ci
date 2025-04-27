package github

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

// TODO: update readme to include the new API client usage
type Client struct {
	Client *api.RESTClient
}

const (
	defaultConcurrency  = 10
	defaultTimeout      = 10 * time.Second
	workflowsPerPage    = 20
	workflowRunsPerPage = 20
	jobsPerPage         = 10
)

type concurrentResult struct {
	Value interface{}
	Error error
}

func runConcurrent(concurrency int, items []interface{}, fn func(item interface{}) (interface{}, error)) []concurrentResult {
	results := make([]concurrentResult, len(items))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i, item := range items {
		wg.Add(1)
		index := i
		currentItem := item

		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result, err := fn(currentItem)
			results[index] = concurrentResult{
				Value: result,
				Error: err,
			}
		}()
	}

	wg.Wait()
	return results
}

func NewClient() (*Client, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}
	return &Client{
		Client: client,
	}, nil
}

// FetchRepositoriesWithWorkflows fetches repositories that have GitHub Actions workflows
func (c *Client) FetchRepositoriesWithWorkflows(names []string) ([]*Repository, error) {
	repos, err := c.fetchRepositories(names)
	if err != nil {
		return nil, fmt.Errorf("error fetching repositories: %w", err)
	}

	if len(repos) == 0 {
		return nil, nil
	}

	// Convert to interface slice for the generic function
	repoItems := make([]interface{}, len(repos))
	for i, repo := range repos {
		repoItems[i] = repo
	}

	results := runConcurrent(defaultConcurrency, repoItems, func(item any) (any, error) {
		repo := item.(*Repository)
		owner, repoName := parseFullName(repo.FullName)
		return c.FetchWorkflowsWithRuns(owner, repoName)
	})

	// Process results
	var successfulRepos []*Repository
	for _, res := range results {
		if res.Error != nil {
			log.Printf("Error fetching workflows: %v", res.Error)
			continue
		}
		successfulRepos = append(successfulRepos, res.Value.(*Repository))
	}

	// Sort repositories by update time (most recent first)
	sort.Slice(successfulRepos, func(i, j int) bool {
		return successfulRepos[i].UpdatedAt.After(successfulRepos[j].UpdatedAt)
	})

	return successfulRepos, nil
}

// FetchWorkflowsWithRuns fetches workflows and their recent runs for a repository
func (c *Client) FetchWorkflowsWithRuns(owner, repo string) (*Repository, error) {
	// Fetch repository info first
	requestUrlRepo := fmt.Sprintf("repos/%s/%s", owner, repo)
	var repository Repository
	err := c.Client.Get(requestUrlRepo, &repository)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository %s/%s: %w", owner, repo, err)
	}

	// Fetch workflows for the repository
	requestUrl := fmt.Sprintf("repos/%s/%s/actions/workflows", owner, repo)
	var workflowsResponse struct {
		Workflows []Workflow `json:"workflows"`
	}

	err = c.Client.Get(requestUrl, &workflowsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflows for %s/%s: %w", owner, repo, err)
	}

	// Convert to interface slice
	workflowItems := make([]interface{}, len(workflowsResponse.Workflows))
	for i, workflow := range workflowsResponse.Workflows {
		workflowItems[i] = workflow
	}

	results := runConcurrent(defaultConcurrency, workflowItems, func(item interface{}) (interface{}, error) {
		workflow := item.(Workflow)

		// Fetch runs for this workflow
		runsUrl := fmt.Sprintf("repos/%s/%s/actions/workflows/%d/runs?per_page=%d",
			owner, repo, workflow.ID, workflowRunsPerPage)

		var runsResponse struct {
			TotalCount   int            `json:"total_count"`
			WorkflowRuns []*WorkflowRun `json:"workflow_runs"`
		}

		err := c.Client.Get(runsUrl, &runsResponse)
		if err != nil {
			workflow.Error = err
			return &workflow, nil // Return with error set but don't fail
		}

		// Fetch jobs for each workflow run
		if len(runsResponse.WorkflowRuns) > 0 {
			workflow.Runs = c.fetchJobsForRuns(owner, repo, runsResponse.WorkflowRuns)
		}

		return &workflow, nil
	})

	// Process results
	var workflows []*Workflow
	for _, res := range results {
		workflows = append(workflows, res.Value.(*Workflow))
	}

	repository.Workflows = workflows
	return &repository, nil
}

// fetchJobsForRuns fetches jobs for a list of workflow runs concurrently
func (c *Client) fetchJobsForRuns(owner, repo string, runs []*WorkflowRun) []*WorkflowRun {
	// Convert to interface slice
	runItems := make([]interface{}, len(runs))
	for i, run := range runs {
		runItems[i] = run
	}

	results := runConcurrent(defaultConcurrency, runItems, func(item interface{}) (interface{}, error) {
		run := item.(*WorkflowRun)

		// Fetch jobs for this run
		jobsUrl := fmt.Sprintf("repos/%s/%s/actions/runs/%d/jobs?per_page=%d",
			owner, repo, run.ID, jobsPerPage)

		var jobsResponse struct {
			TotalCount int    `json:"total_count"`
			Jobs       []*Job `json:"jobs"`
		}

		err := c.Client.Get(jobsUrl, &jobsResponse)
		if err == nil {
			run.Jobs = jobsResponse.Jobs
		}

		return run, nil // Always return run, even if error occurred
	})

	// Convert back to WorkflowRun slice, maintaining the original order
	runsWithJobs := make([]*WorkflowRun, len(runs))
	for i, res := range results {
		runsWithJobs[i] = res.Value.(*WorkflowRun)
	}

	return runsWithJobs
}

// fetchRepositories retrieves repository information for a list of repository names
func (c *Client) fetchRepositories(names []string) ([]*Repository, error) {
	if len(names) == 0 {
		return nil, errors.New("no repository names provided")
	}

	// Convert to interface slice
	nameItems := make([]interface{}, len(names))
	for i, name := range names {
		nameItems[i] = name
	}

	results := runConcurrent(defaultConcurrency, nameItems, func(item interface{}) (interface{}, error) {
		repoName := item.(string)

		repoParts := strings.Split(repoName, "/")
		if len(repoParts) != 2 {
			return nil, fmt.Errorf("invalid repository format: %s (expected 'owner/repo')", repoName)
		}

		requestUrl := fmt.Sprintf("repos/%s/%s", repoParts[0], repoParts[1])
		var response Repository

		err := c.Client.Get(requestUrl, &response)
		if err != nil {
			return nil, err
		}

		return &response, nil
	})

	// Process results
	var repos []*Repository
	for _, res := range results {
		if res.Error != nil {
			log.Printf("Error fetching %v: %v", res.Value, res.Error)
			continue
		}
		repos = append(repos, res.Value.(*Repository))
	}

	return repos, nil
}

// parseFullName splits a full repository name into owner and repo parts
func parseFullName(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", ""
}
