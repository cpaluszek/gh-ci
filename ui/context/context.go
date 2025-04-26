package context

import (
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
)

type ViewType string

const (
	RepoView     ViewType = "repo"
	WorkflowView ViewType = "workflow"
	RunView      ViewType = "run"
)

type Context struct {
	Config            *config.Config
	Client            *github.Client
	Error             error
	ScreenWidth       int
	ScreenHeight      int
	MainContentWidth  int
	MainContentHeight int
	View              ViewType
}
