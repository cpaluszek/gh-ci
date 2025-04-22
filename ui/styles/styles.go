package styles

import "github.com/charmbracelet/lipgloss"

// Color definitions
const (
	// ANSI colors
	BlackANSI        = lipgloss.Color("0")
	RedANSI          = lipgloss.Color("1")
	GreenANSI        = lipgloss.Color("2")
	YellowANSI       = lipgloss.Color("3")
	BlueANSI         = lipgloss.Color("4")
	MagentaANSI      = lipgloss.Color("5")
	CyanANSI         = lipgloss.Color("6")
	WhiteANSI        = lipgloss.Color("7")
	GrayANSI         = lipgloss.Color("8")
	LightGrayANSI    = lipgloss.Color("7")
	LightRedANSI     = lipgloss.Color("9")
	LightGreenANSI   = lipgloss.Color("10")
	LightYellowANSI  = lipgloss.Color("11")
	LightBlueANSI    = lipgloss.Color("12")
	LightMagentaANSI = lipgloss.Color("13")
	LightCyanANSI    = lipgloss.Color("14")
	LightWhiteANSI   = lipgloss.Color("15")
)

// Common styles
var (
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(MagentaANSI)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.NormalBorder()).
			BorderForeground(WhiteANSI).
			BorderBottom(true).
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(RedANSI).
			Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Background(BlackANSI).
			Padding(0, 1).
			Bold(false)

	ActiveWorkflowStyle = lipgloss.NewStyle().
				Foreground(GreenANSI)

	DisabledWorkflowStyle = lipgloss.NewStyle().
				Foreground(LightGrayANSI)

	RowStyle = lipgloss.NewStyle()

	TableHeaderStyle = RowStyle.
				Bold(true)

	SelectedRowStyle = RowStyle.
				Background(BlackANSI)

	SectionContainerStyle = lipgloss.NewStyle().Padding(0, 1)

	SuccessStyle    = lipgloss.NewStyle().Foreground(GreenANSI)
	FailureStyle    = lipgloss.NewStyle().Foreground(RedANSI)
	CanceledStyle   = lipgloss.NewStyle().Foreground(YellowANSI)
	InProgressStyle = lipgloss.NewStyle().Foreground(BlueANSI)
)
