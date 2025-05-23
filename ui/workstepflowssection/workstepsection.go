package workstepflowssection

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
	"github.com/cpaluszek/gh-ci/ui/context"
	"github.com/cpaluszek/gh-ci/ui/keys"
	"github.com/cpaluszek/gh-ci/ui/section"
)

type LogEntry struct {
	Timestamp string
	Level     string
	Message   string
}

type Step struct {
	Number    int
	Title     string
	Status    string
	Duration  string
	Logs      []LogEntry
	Collapsed bool
}

type Model struct {
	section.BaseModel
	steps        []github.Steplog
	cursor       int
	tmpCursor    int
	expandedStep int
	viewport     viewport.Model
	logViewport  viewport.Model
	Job          *github.Job
	headerHeight int
	logHeight    int
	inLogMode    bool
	ready        bool
}

var (
	borderColor   = lipgloss.Color("#404040")
	activeColor   = lipgloss.Color("#61dafb")
	expandedColor = lipgloss.Color("#98fb98")
	mutedColor    = lipgloss.Color("#666666")

	successColor = lipgloss.Color("#98fb98")
	errorColor   = lipgloss.Color("#ff6b6b")
	warningColor = lipgloss.Color("#ffd93d")
	infoColor    = lipgloss.Color("#61dafb")
	runningColor = lipgloss.Color("#ffd93d")

	baseStyle = lipgloss.NewStyle()

	headerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(100)

	stepStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Margin(0, 0, 1, 0).
			Padding(0)

	stepHeaderStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(borderColor)

	stepContentStyle = lipgloss.NewStyle().
				Padding(1, 2)

	logEntryStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 2)
)

func NewModel(ctx *context.Context) Model {
	base := section.NewModel(
		ctx,
		"Logs",
		[]table.Column{
			{
				Title: "",
				Width: 4,
				Grow:  false,
			},
			{
				Title: "Logs",
				Width: 20,
				Grow:  false,
			},
		},
	)

	headerHeight := 2

	logHeight := 15

	vp := viewport.New(
		ctx.MainContentWidth,
		ctx.MainContentHeight-headerHeight,
	)

	logVp := viewport.New(
		ctx.MainContentWidth-2,
		logHeight,
	)

	m := Model{
		BaseModel:    base,
		cursor:       0,
		expandedStep: -1,
		viewport:     vp,
		logViewport:  logVp,
		headerHeight: headerHeight,
		logHeight:    logHeight,
		inLogMode:    false,
		ready:        true,
	}

	m.viewport.GotoTop()
	m.updateViewportContent()
	return m
}

func (m *Model) updateLogger() error {
	info, err := github.ParseGitHubURL(m.Job.GetURL())
	if err != nil {
		return err
	}

	m.steps, err = github.GetLogs(info.User, info.Repo, info.RunID, "1", m.Job.Name)
	return err
}

func (m *Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.GotostepMsg:
		m.Job = msg.RunWithJobs
		m.Table.SetRows(m.BuildRows())
		m.Table.FirstItem()
		cmds = append(cmds, commands.SectionChanged)
		err := m.updateLogger()
		if err != nil {
			return m, tea.Batch(cmds...)
		}
		m.updateViewportContent()

	case tea.WindowSizeMsg:
		headerHeight := m.headerHeight
		footerHeight := 3
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.logViewport = viewport.New(msg.Width-6, m.logHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.logViewport.Width = msg.Width - 6
		}
		m.updateViewportContent()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Keys.Refresh):
			m.updateLogger()
			m.updateViewportContent()

		case key.Matches(msg, keys.Keys.Select):
			if m.cursor == -1 {
				m.cursor = 0
			} else {
				if m.expandedStep == m.cursor {

					m.expandedStep = -1
					m.inLogMode = false
					m.cursor = m.tmpCursor
				} else {

					m.tmpCursor = m.cursor
					m.expandedStep = m.cursor
					m.inLogMode = true
					m.updateLogViewportContent()

					cmds = append(cmds, m.scrollToExpandedStep())
				}
			}
			m.updateViewportContent()

		case key.Matches(msg, keys.Keys.Tab):

			if m.expandedStep != -1 {
				m.expandedStep = -1
				m.inLogMode = false
				m.updateViewportContent()
			}
			m.NextRow()
			m.updateViewportContent()

		case key.Matches(msg, keys.Keys.ShiftTab):

			if m.expandedStep != -1 {
				m.expandedStep = -1
				m.inLogMode = false
				m.updateViewportContent()
			}
			m.PrevRow()
			m.updateViewportContent()

		case msg.Type == tea.KeyUp:
			if m.inLogMode {

				m.logViewport, cmd = m.logViewport.Update(msg)
				m.updateLogViewportContent()
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				return m, tea.Batch(cmds...)
			} else {

				if m.expandedStep != -1 {
					m.expandedStep = -1
					m.inLogMode = false
					m.updateViewportContent()
				}
				m.PrevRow()
				m.updateViewportContent()
			}

		case msg.Type == tea.KeyDown:
			if m.inLogMode {

				m.logViewport, cmd = m.logViewport.Update(msg)
				if cmd != nil {
					cmds = append(cmds, cmd)
				}
				return m, tea.Batch(cmds...)
			} else {

				if m.expandedStep != -1 {
					m.expandedStep = -1
					m.inLogMode = false
					m.updateViewportContent()
				}
				m.NextRow()
				m.updateViewportContent()
			}
		}
	}

	if !m.inLogMode {
		m.viewport, cmd = m.viewport.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else {

		m.logViewport, cmd = m.logViewport.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	table, cmd := m.Table.Update(msg)
	m.Table = table
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) BuildRows() []table.Row {
	return nil
}

func (m Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}

	header := m.renderHeaderView()

	content := m.viewport.View()

	return fmt.Sprintf("%s\n%s", header, content)
}

func (m Model) renderHeaderView() string {
	title := m.Ctx.Styles.Header.Render(fmt.Sprintf(" Workflow Steps%s", strings.Repeat(" ", m.Ctx.MainContentWidth-lipgloss.Width(" Workflow Steps"))))
	return title
}

func (m *Model) updateLogViewportContent() {
	if m.expandedStep >= 0 && m.expandedStep < len(m.steps) {
		step := m.steps[m.expandedStep]
		var logsContent strings.Builder

		for _, logEntry := range step.Logs {
			levelStyle := lipgloss.NewStyle()
			switch strings.ToUpper(logEntry.Level) {
			case "SUCCESS":
				levelStyle = levelStyle.Foreground(successColor)
			case "ERROR":
				levelStyle = levelStyle.Foreground(errorColor)
			case "WARNING":
				levelStyle = levelStyle.Foreground(warningColor)
			case "INFO":
				levelStyle = levelStyle.Foreground(infoColor)
			}

			levelPart := levelStyle.Render("[" + logEntry.Level + "]")
			logLine := fmt.Sprintf("%s %s", levelPart, logEntry.Message)
			logsContent.WriteString(logLine + "\n")
		}

		m.logViewport.SetContent(strings.TrimSuffix(logsContent.String(), "\n"))
	}
}

func (m *Model) updateViewportContent() {
	var content strings.Builder
	totalContentHeight := 0

	expandedStepRenderedStart := 0
	expandedStepRenderedHeight := 0

	if len(m.steps) == 0 {
		noDataText := "No Data"
		noDataStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6b6b")).
			Align(lipgloss.Center).
			Width(m.Ctx.MainContentWidth).
			Height(m.viewport.Height).
			AlignVertical(lipgloss.Center)

		content.WriteString(noDataStyle.Render(noDataText))
		m.viewport.SetContent(content.String())
		return
	}

	for i, step := range m.steps {

		currentStepStyle := stepStyle
		headerStyle := stepHeaderStyle

		if i == m.cursor {
			currentStepStyle = currentStepStyle.BorderForeground(activeColor)
			if i == m.expandedStep {
				currentStepStyle = currentStepStyle.BorderForeground(expandedColor)
				headerStyle = headerStyle.Background(lipgloss.Color("#2d4a2d"))
			}
		}

		stepHeader := m.renderStepHeader(step, i == m.cursor, i == m.expandedStep)

		headerRendered := headerStyle.Width(m.Ctx.MainContentWidth - 4).Render(stepHeader)
		headerHeight := lipgloss.Height(headerRendered)

		stepContent := ""
		stepContentHeight := 0

		if i == m.expandedStep {

			m.updateLogViewportContent()
			stepContent = m.renderStepLogsWithViewport(step)
			stepContentHeight = lipgloss.Height(stepContentStyle.Width(m.Ctx.MainContentWidth - 4).Render(stepContent))
		}

		fullStepHeight := headerHeight + lipgloss.Height(stepStyle.Render(""))
		if stepContent != "" {
			fullStepHeight += stepContentHeight + lipgloss.Height(stepContentStyle.Render(""))
		}

		if i == m.expandedStep {
			expandedStepRenderedStart = totalContentHeight
			expandedStepRenderedHeight = fullStepHeight
		}

		if stepContent != "" {
			fullStep := currentStepStyle.Render(
				lipgloss.JoinVertical(lipgloss.Left,
					headerRendered,
					stepContentStyle.Width(m.Ctx.MainContentWidth-4).Render(stepContent),
				),
			)
			content.WriteString(fullStep + "\n")
		} else {
			fullStep := currentStepStyle.Render(
				headerRendered,
			)
			content.WriteString(fullStep + "\n")
		}
		totalContentHeight += fullStepHeight
	}

	m.viewport.SetContent(content.String())

	if m.expandedStep != -1 {

		if expandedStepRenderedStart+expandedStepRenderedHeight > m.viewport.YOffset+m.viewport.Height {
			m.viewport.SetYOffset(expandedStepRenderedStart + expandedStepRenderedHeight - m.viewport.Height)
		}

		if expandedStepRenderedStart < m.viewport.YOffset {
			m.viewport.SetYOffset(expandedStepRenderedStart)
		}
	}
}

func (m Model) renderStepLogsWithViewport(step github.Steplog) string {

	viewportStyle := lipgloss.NewStyle().
		BorderForeground(borderColor).
		Padding(0)

	return viewportStyle.Render(m.logViewport.View())
}

func (m Model) renderStepHeader(step github.Steplog, isActive, isExpanded bool) string {

	numberStyle := lipgloss.NewStyle().
		Background(borderColor).
		Padding(0, 1).
		Bold(true)

	if isActive {
		numberStyle = numberStyle.Background(activeColor).Foreground(lipgloss.Color("#000000"))
	}
	if isExpanded {
		numberStyle = numberStyle.Background(expandedColor).Foreground(lipgloss.Color("#000000"))
	}

	stepNumber := numberStyle.Render(fmt.Sprintf("%d", step.Number))

	stepTitle := lipgloss.NewStyle().
		Bold(true).
		Render(step.Title)

	statusStyle := lipgloss.NewStyle().Padding(0, 1)
	statusText := strings.ToUpper(step.Status)

	switch step.Status {
	case "success":
		statusStyle = statusStyle.Background(lipgloss.Color("#2d4a2d")).Foreground(successColor)
		statusText = "SUCCESS"
	case "failed":
		statusStyle = statusStyle.Background(lipgloss.Color("#4a2d2d")).Foreground(errorColor)
		statusText = "FAILLED"
	case "running":
		statusStyle = statusStyle.Background(lipgloss.Color("#4a4a2d")).Foreground(warningColor)
		statusText = "RUNNING"
	case "pending":
		statusStyle = statusStyle.Background(lipgloss.Color("#2d3a4a")).Foreground(infoColor)
		statusText = "PENDING"
	case "skipped":
		statusStyle = statusStyle.Background(lipgloss.Color("#4a2d4a")).Foreground(mutedColor)
		statusText = "SKIPPED"
	}

	stepStatus := statusStyle.Render(statusText)

	leftPart := lipgloss.JoinHorizontal(lipgloss.Left,
		stepNumber, " ", stepTitle,
	)

	return lipgloss.JoinHorizontal(lipgloss.Left,
		leftPart,
		strings.Repeat(" ", max(1, m.Ctx.MainContentWidth-8-lipgloss.Width(leftPart)-lipgloss.Width(stepStatus))),
		stepStatus,
	)
}

func (m *Model) NumRows() int {
	return len(m.steps)
}

func (m *Model) NextRow() int {
	if m.cursor < len(m.steps)-1 {
		m.cursor++
	}
	return m.cursor
}

func (m *Model) PrevRow() int {
	if m.cursor > 0 {
		m.cursor--
	}
	return m.cursor
}

func (m *Model) GetCurrentRow() github.RowData {
	if m.cursor >= 0 && m.cursor < len(m.steps) {
		return nil
	}
	return nil
}

func (m *Model) SetIsLoading(val bool) {
	m.IsLoading = val
}

func (m *Model) UpdateContext(ctx *context.Context) {
	m.Ctx = ctx

	m.viewport.Width = m.Ctx.MainContentWidth
	m.viewport.Height = m.Ctx.MainContentHeight - m.headerHeight

	m.logViewport.Width = m.Ctx.MainContentWidth - 6

	m.updateViewportContent()
}

func (m *Model) Fetch() []tea.Cmd {
	return nil
}

func (m *Model) IsInLogMode() bool {
	return m.inLogMode
}

func (m *Model) scrollToExpandedStep() tea.Cmd {
	return func() tea.Msg {

		return nil
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
