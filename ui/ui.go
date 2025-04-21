package ui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/footer"
	"github.com/cpaluszek/pipeye/ui/components/reposection"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/section"
)

type Model struct {
	footer  footer.Model
	ctx     context.Context
	spinner spinner.Model
	loading bool
	repos   section.Section
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	m := Model{
		footer:  footer.NewModel(),
		spinner: s,
		loading: true,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commands.InitConfig,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case commands.ConfigInitMsg:
		m.ctx.Config = msg.Config
		s := reposection.NewModel(m.ctx)
		m.repos = &s
		return m, commands.InitClient(m.ctx.Config.Github.Token)

	case commands.ClientInitMsg:
		m.ctx.Client = msg.Client
		m.loading = false
		return m, commands.FetchRepositories(m.ctx.Client)

	case commands.ErrorMsg:
		log.Println("Error:", msg.Error)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	}

	if !m.loading {
		var sectionCmd tea.Cmd
		m.repos, sectionCmd = m.repos.Update(msg)

		cmds = append(
			cmds,
			sectionCmd,
		)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}

	if m.loading {
		s.WriteString(m.spinner.View() + "\n")
	} else {
		s.WriteString(m.repos.View())
	}

	s.WriteString(m.footer.View())
	return s.String()
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
}
