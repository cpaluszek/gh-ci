package context

import (
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
)

type Context struct {
	Config            *config.Config
	Client            *github.Client
	Error             error
	ScreenWidth       int
	ScreenHeight      int
	MainContentWidth  int
	MainContentHeight int
}
