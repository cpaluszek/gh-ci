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

// TODO: fix loading
// TODO: fix table scrolling while in log mode
type Model struct {
	section.BaseModel
	steps        []github.Steplog
	expandedStep int
	logViewport  viewport.Model
	Job          *github.Job
	logHeight    int
	inLogMode    bool
	error        string
}

var logHeight = 15

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Workflow Steps",
		[]table.Column{
			{
				Title: "#",
				Width: 4,
				Grow:  false,
			},
			{
				Title: "Title",
				Width: 40,
				Grow:  true,
			},
			{
				Title: "Status",
				Width: 15,
				Grow:  false,
			},
			{
				Title: "Duration",
				Width: 15,
				Grow:  false,
			},
		},
	)

	logVp := viewport.New(
		ctx.MainContentWidth-2,
		logHeight,
	)

	m := Model{
		BaseModel:    base,
		expandedStep: -1,
		logViewport:  logVp,
		logHeight:    logHeight,
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
		err := m.getLogs()
		if err != nil {
			m.error = err.Error()
		} else {
			m.Table.SetRows(m.BuildRows())
			m.Table.FirstItem()
		}
		m.SetIsLoading(false)
		cmds = append(cmds, commands.SectionChanged)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Refresh):
			err := m.getLogs()
			if err != nil {
				m.error = err.Error()
			} else {
				m.Table.SetRows(m.BuildRows())
			}

		case key.Matches(msg, keys.Keys.Select):
			currentIndex := m.Table.GetCurrItem()
			if currentIndex >= 0 && currentIndex < len(m.steps) {
				if m.expandedStep == currentIndex {
					m.expandedStep = -1
					m.inLogMode = false
					m.Table.SetDimensions(m.GetDimensions())
					m.Table.SyncViewPortContent()
				} else {
					m.expandedStep = currentIndex
					m.inLogMode = true
					m.updateLogViewportContent()
					tableDim := constants.Dimensions{
						Width:  m.Ctx.MainContentWidth,
						Height: m.Ctx.MainContentHeight - logHeight - 2,
					}
					m.Table.SetDimensions(tableDim)
					m.Table.SyncViewPortContent()
				}
			}

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

	if m.inLogMode && cmd == nil {
		m.logViewport, cmd = m.logViewport.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	table, cmd := m.Table.Update(msg)
	m.Table = table
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) getLogs() error {
	info, err := github.ParseGitHubURL(m.Job.GetURL())
	if err != nil {
		return err
	}

	m.steps, err = m.Ctx.Client.GetLogs(info.User, info.Repo, info.RunID, "1", m.Job.Name)
	return err
}

func (m Model) BuildRows() []table.Row {
	var rows []table.Row

	for _, step := range m.steps {
		statusText := strings.ToUpper(step.Status)
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
			fmt.Sprintf("%d", step.Number),
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

	tableView := m.Table.View()

	if m.expandedStep >= 0 && m.expandedStep < len(m.steps) {
		logsView := m.renderExpandedStepLogs()
		return lipgloss.JoinVertical(
			lipgloss.Left,
			tableView,
			logsView,
		)
	}

	return tableView
}

func (m Model) renderExpandedStepLogs() string {
	if m.expandedStep < 0 || m.expandedStep >= len(m.steps) {
		return ""
	}

	step := m.steps[m.expandedStep]
	title := fmt.Sprintf(" Logs for Step %d: %s", step.Number, step.Title)

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

	m.logViewport.Width = m.Ctx.MainContentWidth - 6
}

func (m *Model) Fetch() []tea.Cmd {
	return nil
}

func (m *Model) GetCurrentRow() github.RowData {
	return nil
}

func (m *Model) IsInLogMode() bool {
	return m.inLogMode
}
