package ui

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/github"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/footer"
	"github.com/cpaluszek/pipeye/ui/components/reposection"
	"github.com/cpaluszek/pipeye/ui/components/sidebar"
	"github.com/cpaluszek/pipeye/ui/constants"
	"github.com/cpaluszek/pipeye/ui/context"
	"github.com/cpaluszek/pipeye/ui/runsection"
	"github.com/cpaluszek/pipeye/ui/section"
	"github.com/cpaluszek/pipeye/ui/workflowssection"
)

type Model struct {
	footer   footer.Model
	ctx      *context.Context
	repos    section.Section
	worflows section.Section
	run      section.Section
	sidebar  sidebar.Model
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
	r := runsection.NewModel(m.ctx)
	m.run = &r
	sidebar := sidebar.NewModel(m.ctx)
	m.sidebar = sidebar

	return m
}

func (m Model) Init() tea.Cmd {
	m.ctx.View = context.RepoView
	return commands.InitConfig
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.footer, cmd = m.footer.Update(msg)
			return m, cmd
		case "j", "down":
			m.GetCurrentSection().NextRow()
			m.OnSelectedRowChanged()
		case "k", "up":
			m.GetCurrentSection().PrevRow()
			m.OnSelectedRowChanged()
		case "enter":
			switch m.ctx.View {
			case context.RepoView:
				repo := m.repos.GetCurrentRow()
				m.ctx.View = context.WorkflowView
				return m, commands.GoToWorkflow(repo)

			case context.WorkflowView:
				workflowRun := m.worflows.GetCurrentRow()
				m.ctx.View = context.RunView
				return m, commands.GoToRun(workflowRun)
			}
		case "esc", "backspace":
			switch m.ctx.View {
			case context.WorkflowView:
				m.ctx.View = context.RepoView
				m.OnSelectedRowChanged()
			case context.RunView:
				m.ctx.View = context.WorkflowView
				m.OnSelectedRowChanged()
			}

		}
	case commands.ConfigInitMsg:
		m.ctx.Config = msg.Config
		return m, commands.InitClient(m.ctx.Config.Github.Token)

	case commands.ClientInitMsg:
		m.ctx.Client = msg.Client
		cmds = append(cmds, m.repos.Fetch()...)

	case commands.SectionChangedMsg:
		m.OnSelectedRowChanged()

	case commands.ErrorMsg:
		log.Println("Error:", msg.Error)
		return m, nil

	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	}

	sectionCmd := m.updateCurrentSection(msg)
	m.sidebar.UpdateProgramContext(m.ctx)

	var footerCmd tea.Cmd
	m.footer, footerCmd = m.footer.Update(msg)

	cmds = append(
		cmds,
		sectionCmd,
		footerCmd,
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}

	s.WriteString("\n")
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.GetCurrentSection().View(),
		m.sidebar.View(),
	)

	s.WriteString(content)
	s.WriteString("\n")

	s.WriteString(m.footer.View())
	return s.String()
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	footerHeight := 1
	headerHeight := 1
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.ctx.MainContentWidth = msg.Width - constants.SideBarWidth
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
	case context.RunView:
		m.run.UpdateContext(m.ctx)
		m.run, cmd = m.run.Update(msg)
	}

	return cmd
}

func (m *Model) GetCurrentSection() section.Section {
	switch m.ctx.View {
	case context.RepoView:
		return m.repos
	case context.WorkflowView:
		return m.worflows
	case context.RunView:
		return m.run
	}
	return m.repos
}

func (m *Model) OnSelectedRowChanged() {
	// Sidebar sync
	currentRow := m.GetCurrentSection().GetCurrentRow()
	if currentRow == nil {
		m.sidebar.SetContent("")
		return
	}

	switch m.ctx.View {
	case context.RepoView:
		if repo, ok := currentRow.(*github.Repository); ok {
			m.sidebar.GenerateRepoSidebarContent(repo)
		}
	case context.WorkflowView:
		if workflowRun, ok := currentRow.(*github.WorkflowRun); ok {
			m.sidebar.GenerateWorkflowSidebarContent(workflowRun)
		}
	case context.RunView:
		if jobData, ok := currentRow.(*github.Job); ok {
			m.sidebar.GenerateRunSidebarContent(jobData)
		}
	}
}
