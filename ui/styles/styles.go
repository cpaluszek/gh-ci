package styles

import "github.com/charmbracelet/lipgloss"

var (
	PrimaryBorder   = lipgloss.AdaptiveColor{Light: "013", Dark: "008"}
	SecondaryBorder = lipgloss.AdaptiveColor{Light: "008", Dark: "007"}

	SelectedBackground = lipgloss.AdaptiveColor{Light: "006", Dark: "008"}

	PrimaryText   = lipgloss.AdaptiveColor{Light: "000", Dark: "015"}
	SecondaryText = lipgloss.AdaptiveColor{Light: "244", Dark: "251"}
	FaintText     = lipgloss.AdaptiveColor{Light: "007", Dark: "254"}

	SuccessText = lipgloss.AdaptiveColor{Light: "002", Dark: "002"}
	WarningText = lipgloss.AdaptiveColor{Light: "003", Dark: "003"}
	ErrorText   = lipgloss.AdaptiveColor{Light: "001", Dark: "001"}

	SuccessColor    = SuccessText
	WarningColor    = WarningText
	ErrorColor      = ErrorText
	InProgressColor = lipgloss.AdaptiveColor{Light: "004", Dark: "004"}
	SkippedColor    = FaintText
)

// Common styles
var (
	SpinnerStyle = lipgloss.NewStyle().Bold(true).Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.NormalBorder()).
			BorderForeground(PrimaryBorder).
			BorderBottom(true).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)

	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryText).
			Bold(true)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(ErrorText).
			Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Background(SelectedBackground).
			Padding(0, 1)

	RowStyle = lipgloss.NewStyle()

	TableHeaderStyle = RowStyle.
				Bold(true)

	SelectedRowStyle = RowStyle.Background(SelectedBackground)

	SectionContainerStyle = lipgloss.NewStyle().Padding(0, 1)

	SideBarStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.NormalBorder()).
			BorderForeground(PrimaryBorder).
			BorderBottom(false).
			BorderTop(false).
			BorderLeft(true).
			BorderRight(false)

	SuccessStyle    = lipgloss.NewStyle().Foreground(SuccessColor)
	FailureStyle    = lipgloss.NewStyle().Foreground(ErrorColor)
	CanceledStyle   = lipgloss.NewStyle().Foreground(WarningColor)
	InProgressStyle = lipgloss.NewStyle().Foreground(InProgressColor)
	SkippedStyle    = lipgloss.NewStyle().Foreground(SkippedColor)
	DefaultStyle    = lipgloss.NewStyle().Foreground(PrimaryText)
)

// TODO: add plain text fallback
// Nerd Font workflow status symbols
const (
	// Status symbols
	SuccessSymbol    = "󰄬 "
	FailureSymbol    = "󰅚 "
	CanceledSymbol   = "󰔛 "
	SkippedSymbol    = "󰒭 "
	NeutralSymbol    = "󰘿 "
	InProgressSymbol = "󰑮 "
	QueuedSymbol     = "󰥔 "

	// Job status dot variants
	JobSuccessDot    = "󰄯 "
	JobFailureDot    = "󰅙 "
	JobCanceledDot   = "󰅚 "
	JobSkippedDot    = "○ "
	JobInProgressDot = "◌ "
)

func GetStatusSymbol(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return SuccessStyle.Render(SuccessSymbol)
		case "failure", "timed_out", "startup_failure":
			return FailureStyle.Render(FailureSymbol)
		case "cancelled":
			return CanceledStyle.Render(CanceledSymbol)
		case "skipped":
			return SkippedStyle.Render(SkippedSymbol)
		case "neutral":
			return DefaultStyle.Render(NeutralSymbol)
		default:
			return DefaultStyle.Render(NeutralSymbol)
		}
	case "in_progress":
		return InProgressStyle.Render(InProgressSymbol)
	case "queued":
		return InProgressStyle.Render(QueuedSymbol)
	default:
		return DefaultStyle.Render(NeutralSymbol)
	}
}

func GetJobStatusSymbol(conclusion string) string {
	switch conclusion {
	case "success":
		return SuccessStyle.Render(JobSuccessDot)
	case "failure", "timed_out", "startup_failure":
		return FailureStyle.Render(JobFailureDot)
	case "cancelled":
		return CanceledStyle.Render(JobCanceledDot)
	case "skipped":
		return SkippedStyle.Render(JobSkippedDot)
	case "in_progress":
		return DefaultStyle.Render(JobInProgressDot)
	default:
		return DefaultStyle.Render(JobSkippedDot)
	}
}
