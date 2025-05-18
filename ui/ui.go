package ui

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/gh-ci/config"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/commands"
	"github.com/cpaluszek/gh-ci/ui/components/footer"
	"github.com/cpaluszek/gh-ci/ui/components/reposection"
	"github.com/cpaluszek/gh-ci/ui/components/sidebar"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
	"github.com/cpaluszek/gh-ci/ui/runsection"
	"github.com/cpaluszek/gh-ci/ui/section"
	"github.com/cpaluszek/gh-ci/ui/styles"
	"github.com/cpaluszek/gh-ci/ui/workflowssection"
)

type Model struct {
	footer   footer.Model
	ctx      *context.Context
	repos    section.Section
	worflows section.Section
	run      section.Section
	sidebar  sidebar.Model
}

func NewModel(cfg *config.Config) Model {
	theme := styles.DefaultTheme
	styles := styles.BuildStyles(*theme)
	m := Model{
		ctx: &context.Context{
			ScreenWidth:  0,
			ScreenHeight: 0,
			Config:       cfg,
			Theme:        theme,
			Styles:       &styles,
		},
	}
	f := footer.NewModel(m.ctx)
	m.footer = f

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
	return commands.InitClient()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Quit):
			m.footer, cmd = m.footer.Update(msg)
			return m, cmd
		case key.Matches(msg, keys.Keys.Down):
			m.GetCurrentSection().NextRow()
			m.OnSelectedRowChanged()
		case key.Matches(msg, keys.Keys.Up):
			m.GetCurrentSection().PrevRow()
			m.OnSelectedRowChanged()
		case key.Matches(msg, keys.Keys.Select):
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
		case key.Matches(msg, keys.Keys.Return):
			switch m.ctx.View {
			case context.WorkflowView:
				m.ctx.View = context.RepoView
				m.OnSelectedRowChanged()
			case context.RunView:
				m.ctx.View = context.WorkflowView
				m.OnSelectedRowChanged()
			}
		case key.Matches(msg, keys.Keys.Help):
			if m.footer.Help.ShowAll {
				m.ctx.MainContentHeight = m.ctx.MainContentHeight + constants.HelpHeight - constants.FooterHeight
			} else {
				m.ctx.MainContentHeight = m.ctx.MainContentHeight + constants.FooterHeight - constants.HelpHeight
			}
		}
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
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.ctx.MainContentWidth = msg.Width - constants.SideBarWidth
	m.ctx.MainContentHeight = msg.Height - constants.FooterHeight - constants.HeaderHeight
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
	return nil
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
