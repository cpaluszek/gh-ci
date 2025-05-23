package context

import (
	"github.com/cpaluszek/gh-ci/config"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/styles"
)

type ViewType string

const (
	RepoView     ViewType = "repo"
	WorkflowView ViewType = "workflow"
	RunView      ViewType = "run"
	LogStepView  ViewType = "step"
	LogView      ViewType = "log"
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
