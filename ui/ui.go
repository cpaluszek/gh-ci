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
	"github.com/cpaluszek/pipeye/ui/workflowssection"
)

type Model struct {
	footer   footer.Model
	ctx      *context.Context
	repos    section.Section
	worflows section.Section
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
	w := workflowssection.NewModel(m.ctx)
	m.worflows = &w

	return m
}

func (m Model) Init() tea.Cmd {
	m.ctx.View = context.RepoView
	return commands.InitConfig
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j", "down":
			m.GetCurrentSection().NextRow()
		case "k", "up":
			m.GetCurrentSection().PrevRow()
		case "enter":
			switch m.ctx.View {
			case context.RepoView:
				repo := m.repos.GetCurrentRow()
				m.ctx.View = context.WorkflowView

				return m, commands.GoToWorkflow(repo)
			}
		case "esc", "backspace":
			switch m.ctx.View {
			case context.WorkflowView:
				m.ctx.View = context.RepoView
			}

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

	sectionCmd := m.updateCurrentSection(msg)

	cmds = append(
		cmds,
		sectionCmd,
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString("\n")
	currentSection := m.GetCurrentSection()
	s.WriteString(currentSection.View())
	s.WriteString("\n")

	s.WriteString(m.footer.View())
	return s.String()
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	footerHeight := 1
	headerHeight := 1
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.ctx.MainContentWidth = msg.Width
	m.ctx.MainContentHeight = msg.Height - footerHeight - headerHeight
	m.footer.SetWidth(msg.Width)
}

func (m *Model) updateCurrentSection(msg tea.Msg) (cmd tea.Cmd) {
	switch m.ctx.View {
	case context.RepoView:
		m.repos.UpdateContext(m.ctx)
		m.repos, cmd = m.repos.Update(msg)
	case context.WorkflowView:
		m.worflows.UpdateContext(m.ctx)
		m.worflows, cmd = m.worflows.Update(msg)
	}

	return cmd
}

func (m *Model) GetCurrentSection() section.Section {
	switch m.ctx.View {
	case context.RepoView:
		return m.repos
	case context.WorkflowView:
		return m.worflows
	}
	return m.repos
}
