package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/ui"
)

type BaseView struct {
	Viewport       viewport.Model
	Client         *github.Client
	Loading        bool
	Error          error
	Spinner        spinner.Model
	StatusBarStyle lipgloss.Style
}

func NewBaseView(vp viewport.Model, client *github.Client) BaseView {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = ui.SpinnerStyle

	return BaseView{
		Viewport:       vp,
		Client:         client,
		Loading:        true,
		Spinner:        s,
		StatusBarStyle: ui.StatusStyle,
	}
}

func (b *BaseView) UpdateSize(width, height int) {
	b.Viewport.Width = width
	b.Viewport.Height = height - 1 // reserve space for status
	b.StatusBarStyle = ui.StatusStyle.Width(width)
}
