package github

import (
	"time"
)

type Repository struct {
	ID             int64       `json:"id"` // Changed from string to int64
	Name           string      `json:"name"`
	FullName       string      `json:"full_name"`
	URL            string      `json:"html_url"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Language       string      `json:"language"`
	IsPrivate      bool        `json:"private"`
	StargazerCount int         `json:"stargazers_count"`
	Workflows      []*Workflow `json:"-"` // Not directly from the API
	Error          error       `json:"-"` // Not from the API
}

// Workflow represents a GitHub Actions workflow
type Workflow struct {
	ID    int64          `json:"id"`
	Name  string         `json:"name"`
	State string         `json:"state"`
	URL   string         `json:"html_url"`
	Runs  []*WorkflowRun `json:"-"` // Not from direct API response
	Error error          `json:"-"` // Not from API
}

// WorkflowRun represents a run of a GitHub Actions workflow
type WorkflowRun struct {
	ID           int64     `json:"id"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DisplayTitle string    `json:"display_title"`
	Event        string    `json:"event"`
	URL          string    `json:"html_url"`
	HeadBranch   string    `json:"head_branch"`
	HeadCommit   Commit    `json:"head_commit"`
	Jobs         []*Job    `json:"-"` // Fetched separately
}

// Commit represents a git commit
type Commit struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

// Job represents a job in a workflow run
type Job struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
	URL         string    `json:"html_url"`
	Steps       []Step    `json:"steps"`
}

// Step represents a step in a workflow job
type Step struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int       `json:"number"`
	Completed   bool      `json:"completed"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type RowData interface {
	GetName() string
	GetURL() string
}

func (r Repository) GetName() string {
	return r.Name
}

func (r Repository) GetURL() string {
	return r.URL
}

func (w WorkflowRun) GetName() string {
	return w.DisplayTitle
}

func (w WorkflowRun) GetURL() string {
	return w.URL
}

func (j Job) GetName() string {
	return j.Name
}

func (j Job) GetURL() string {
	return j.URL
}
