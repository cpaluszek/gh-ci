package github

import (
	"fmt"

	gh "github.com/google/go-github/v71/github"
)

type RepositoryData struct {
	Repository          *gh.Repository
	WorkflowRunWithJobs []*WorkflowWithRuns
	Error               error
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

type RowData interface {
	GetID() string
	GetName() string
	GetURL() string
}

func (r RepositoryData) GetID() string {
	return r.Repository.GetNodeID()
}

func (r RepositoryData) GetName() string {
	return r.Repository.GetFullName()
}

func (r RepositoryData) GetURL() string {
	return r.Repository.GetHTMLURL()
}

func (w WorkflowRunWithJobs) GetID() string {
	if w.Run == nil {
		return ""
	}
	return fmt.Sprintf("%d", w.Run.GetID())
}

func (w WorkflowRunWithJobs) GetName() string {
	if w.Run == nil {
		return ""
	}
	return w.Run.GetDisplayTitle()
}

func (w WorkflowRunWithJobs) GetURL() string {
	if w.Run == nil {
		return ""
	}
	return w.Run.GetHTMLURL()
}
