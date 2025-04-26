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

	PullRequestColor = lipgloss.AdaptiveColor{Light: "004", Dark: "004"}
	PushColor        = lipgloss.AdaptiveColor{Light: "003", Dark: "003"}
	ScheduleColor    = lipgloss.AdaptiveColor{Light: "007", Dark: "007"}
	PlayColor        = lipgloss.AdaptiveColor{Light: "005", Dark: "005"}
	IssueColor       = lipgloss.AdaptiveColor{Light: "001", Dark: "001"}
	DeploymentColor  = lipgloss.AdaptiveColor{Light: "002", Dark: "002"}
	TagColor         = lipgloss.AdaptiveColor{Light: "006", Dark: "006"}
	WebHookColor     = lipgloss.AdaptiveColor{Light: "003", Dark: "003"}
	ForkColor        = lipgloss.AdaptiveColor{Light: "008", Dark: "008"}
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

	PullRequestStyle = lipgloss.NewStyle().Foreground(PullRequestColor).Bold(true)
	PushStyle        = lipgloss.NewStyle().Foreground(PushColor).Bold(true)
	ScheduleStyle    = lipgloss.NewStyle().Foreground(ScheduleColor).Bold(true)
	PlayStyle        = lipgloss.NewStyle().Foreground(PlayColor).Bold(true)
	IssueStyle       = lipgloss.NewStyle().Foreground(IssueColor).Bold(true)
	DeploymentStyle  = lipgloss.NewStyle().Foreground(DeploymentColor).Bold(true)
	TagStyle         = lipgloss.NewStyle().Foreground(TagColor).Bold(true)
	WebHookStyle     = lipgloss.NewStyle().Foreground(WebHookColor).Bold(true)
	ForkStyle        = lipgloss.NewStyle().Foreground(ForkColor).Bold(true)
)
