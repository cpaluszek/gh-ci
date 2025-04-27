package context

import (
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/styles"
)

type ViewType string

const (
	RepoView     ViewType = "repo"
	WorkflowView ViewType = "workflow"
	RunView      ViewType = "run"
)

type Context struct {
	Config            *config.Config
	Theme             *styles.Theme
	Styles            *styles.Styles
	Client            *github.Client
	Error             error
	ScreenWidth       int
	ScreenHeight      int
	MainContentWidth  int
	MainContentHeight int
	View              ViewType
}
