package ui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/footer"
	"github.com/cpaluszek/pipeye/ui/components/reposection"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/section"
)

type Model struct {
	footer footer.Model
	ctx    *context.Context
	repos  section.Section
}

func NewModel() Model {
	m := Model{
		footer: footer.NewModel(),
		ctx: &context.Context{
			ScreenWidth:  0,
			ScreenHeight: 0,
		},
	}
	s := reposection.NewModel(m.ctx)
	m.repos = &s

	return m
}

func (m Model) Init() tea.Cmd {
	return commands.InitConfig
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "j", "down":
			m.repos.NextRow()
		case "k", "up":
			m.repos.PrevRow()
		}
	case commands.ConfigInitMsg:
		m.ctx.Config = msg.Config
		return m, commands.InitClient(m.ctx.Config.Github.Token)

	case commands.ClientInitMsg:
		m.ctx.Client = msg.Client
		cmds = append(cmds, m.repos.Fetch()...)

	case commands.ErrorMsg:
		log.Println("Error:", msg.Error)
		return m, nil

	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	}

	m.repos.UpdateContext(m.ctx)

	sectionCmd := m.updateCurrentSection(msg)

	cmds = append(
		cmds,
		sectionCmd,
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString(m.repos.View())
	s.WriteString("\n")

	s.WriteString(m.footer.View())
	return s.String()
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	footerHeight := 1
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.ctx.MainContentWidth = msg.Width
	m.ctx.MainContentHeight = msg.Height - footerHeight
	m.footer.SetWidth(msg.Width)
}

func (m *Model) updateCurrentSection(msg tea.Msg) (cmd tea.Cmd) {
	// TODO: get current section
	m.repos, cmd = m.repos.Update(msg)

	return cmd
}
