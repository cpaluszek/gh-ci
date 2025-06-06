package stepsection

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/commands"
	"github.com/cpaluszek/gh-ci/ui/components/table"
	"github.com/cpaluszek/gh-ci/ui/constants"
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
	"github.com/cpaluszek/gh-ci/ui/section"
)

type Model struct {
	section.BaseModel
	steps        []github.Steplog
	expandedStep int
	logViewport  viewport.Model
	Job          *github.Job
	inLogMode    bool
	error        string
}

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Workflow Steps",
		[]table.Column{
			{
				Title: "Step",
				Width: 40,
				Grow:  true,
			},
			{
				Title: "Status",
				Width: 18,
				Grow:  false,
			},
			{
				Title: "Duration",
				Width: 12,
				Grow:  false,
			},
		},
	)

	logVp := viewport.New(
		ctx.MainContentWidth,
		ctx.MainContentHeight-2,
	)

	m := Model{
		BaseModel:    base,
		expandedStep: -1,
		logViewport:  logVp,
		inLogMode:    false,
	}

	return m
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.GotostepMsg:
		m.Job = msg.RunWithJobs
		m.expandedStep = -1
		m.inLogMode = false
		return m, tea.Batch(m.Fetch()...)

	case commands.LogsMsg:
		m.steps = msg.Steps
		m.SetIsLoading(false)
		m.Table.SetRows(m.BuildRows())
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Select):
			currentIndex := m.Table.GetCurrItem()
			if !m.inLogMode && currentIndex >= 0 && currentIndex < len(m.steps) {
				m.expandedStep = currentIndex
				m.inLogMode = true
				m.updateLogViewportContent()
			}

		case key.Matches(msg, keys.Keys.Return):
			m.expandedStep = -1
			m.inLogMode = false

		case key.Matches(msg, keys.Keys.OpenGitHub):
            if m.Job == nil {
                return m, nil
            }

            url := m.Job.URL
            if url == "" {
                return m, nil
            }
            return m, commands.OpenBrowser(url)

		case key.Matches(msg, keys.Keys.Up) || key.Matches(msg, keys.Keys.Down):
			if m.inLogMode {
				m.logViewport, cmd = m.logViewport.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
			} else {
				m.expandedStep = -1
				table, cmd := m.Table.Update(msg)
				cmds = append(cmds, cmd)
				m.Table = table
			}
		}
	}

	if m.inLogMode {
		m.logViewport, cmd = m.logViewport.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		table, cmd := m.Table.Update(msg)
		m.Table = table
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) BuildRows() []table.Row {
	var rows []table.Row

	for _, step := range m.steps {
		statusText := step.Status
		statusDisplay := statusText

		switch step.Status {
		case "success":
			statusDisplay = m.Ctx.Styles.Success.Render(statusText)
		case "failed":
			statusDisplay = m.Ctx.Styles.Error.Render(statusText)
		case "running":
			statusDisplay = m.Ctx.Styles.Warning.Render(statusText)
		case "pending":
			statusDisplay = m.Ctx.Styles.Info.Render(statusText)
		case "skipped":
			statusDisplay = m.Ctx.Styles.Skipped.Render(statusText)
		}

		rows = append(rows, table.Row{
			step.Title,
			statusDisplay,
			step.Duration,
		})
	}

	return rows
}

func (m Model) View() string {
	if m.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(m.Ctx.Theme.Colors.Error).
			Align(lipgloss.Center).
			Width(m.Ctx.MainContentWidth).
			Height(m.Ctx.MainContentHeight).
			AlignVertical(lipgloss.Center)

		return errorStyle.Render(m.error)
	}

	if m.expandedStep >= 0 && m.expandedStep < len(m.steps) {
		return m.renderExpandedStepLogs()
	}

	return m.Table.View()
}

func (m Model) renderExpandedStepLogs() string {
	if m.expandedStep < 0 || m.expandedStep >= len(m.steps) {
		return ""
	}

	step := m.steps[m.expandedStep]
	title := fmt.Sprintf(" Step %s:", step.Title)

	logsHeader := m.Ctx.Styles.Header.Render(title)
	logsContent := m.logViewport.View()

	logBox := m.Ctx.Styles.Default.
		Width(m.Ctx.MainContentWidth).
		Render(logsContent)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		logsHeader,
		logBox,
	)
}

func (m *Model) updateLogViewportContent() {
	if m.expandedStep >= 0 && m.expandedStep < len(m.steps) {
		step := m.steps[m.expandedStep]
		var logsContent strings.Builder

		for _, logEntry := range step.Logs {
			levelStyle := lipgloss.NewStyle()

			switch strings.ToUpper(logEntry.Level) {
			case "SUCCESS":
				levelStyle = levelStyle.Foreground(m.Ctx.Theme.Colors.Success)
			case "ERROR":
				levelStyle = levelStyle.Foreground(m.Ctx.Theme.Colors.Error)
			case "WARNING":
				levelStyle = levelStyle.Foreground(m.Ctx.Theme.Colors.Warning)
			case "INFO":
				levelStyle = levelStyle.Foreground(m.Ctx.Theme.Colors.Info)
			}

			levelPart := levelStyle.Render("[" + logEntry.Level + "]")
			logLine := fmt.Sprintf("%s %s", levelPart, logEntry.Message)
			logsContent.WriteString(logLine + "\n")
		}

		m.logViewport.SetContent(strings.TrimSuffix(logsContent.String(), "\n"))
		m.logViewport.GotoTop()
	}
}

func (m *Model) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.Ctx.MainContentWidth,
		Height: m.Ctx.MainContentHeight,
	}
}

func (m *Model) NumRows() int {
	return len(m.steps)
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
	m.Table.SetIsLoading(val)
}

func (m *Model) UpdateContext(ctx *context.Context) {
	m.Ctx = ctx
	m.Table.UpdateContext(ctx)
	m.Table.SetDimensions(m.GetDimensions())
	m.Table.SyncViewPortContent()

	m.logViewport.Width = m.Ctx.MainContentWidth
	m.logViewport.Height = m.Ctx.MainContentHeight - 2
}

func (m *Model) Fetch() []tea.Cmd {
	if m == nil || m.Job == nil {
		return nil
	}

	var cmds []tea.Cmd
	tableCmd := m.Table.StartLoadingSpinner()
	fetchCmd := commands.FetchStepLogs(m.Ctx.Client, m.Job)

	cmds = append(cmds, tableCmd, fetchCmd)
	m.SetIsLoading(true)
	return cmds
}

func (m *Model) GetCurrentRow() github.RowData {
	return nil
}
