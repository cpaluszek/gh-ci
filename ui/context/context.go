package context

import (
	"github.com/cpaluszek/gh-actions/config"
	"github.com/cpaluszek/gh-actions/github"
	"github.com/cpaluszek/gh-actions/ui/styles"
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
