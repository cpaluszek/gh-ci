package github

import (
	"fmt"

	gh "github.com/google/go-github/v71/github"
)

// TODO: clarify types and variables naming
type Repository struct {
	Info      *gh.Repository
	Workflows []*Workflow
	Error     error
}

type Workflow struct {
	Info  *gh.Workflow
	Runs  []*WorkflowRun
	Error error
}

type WorkflowRun struct {
	Info *gh.WorkflowRun
	Jobs []*gh.WorkflowJob
}

type Job struct {
	Job *gh.WorkflowJob
}

type RowData interface {
	GetID() string
	GetName() string
	GetURL() string
}

func (r Repository) GetID() string {
	return r.Info.GetNodeID()
}

func (r Repository) GetName() string {
	return r.Info.GetFullName()
}

func (r Repository) GetURL() string {
	return r.Info.GetHTMLURL()
}

func (w WorkflowRun) GetID() string {
	if w.Info == nil {
		return ""
	}
	return fmt.Sprintf("%d", w.Info.GetID())
}

func (w WorkflowRun) GetName() string {
	if w.Info == nil {
		return ""
	}
	return w.Info.GetDisplayTitle()
}

func (w WorkflowRun) GetURL() string {
	if w.Info == nil {
		return ""
	}
	return w.Info.GetHTMLURL()
}

func (j Job) GetID() string {
	if j.Job == nil {
		return ""
	}
	return j.Job.GetNodeID()
}

func (j Job) GetName() string {
	if j.Job == nil {
		return ""
	}
	return j.Job.GetName()
}

func (j Job) GetURL() string {
	if j.Job == nil {
		return ""
	}
	return j.Job.GetHTMLURL()
}
