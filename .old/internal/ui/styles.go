package ui

import "github.com/charmbracelet/lipgloss"

// type Theme struct {
// 	PrimaryText   lipgloss.AdaptiveColor
// 	SecondaryText lipgloss.AdaptiveColor

// 	ErrorText lipgloss.AdaptiveColor

// 	SelectedBackground lipgloss.AdaptiveColor

// 	SuccessColor lipgloss.AdaptiveColor
// 	FailureColor lipgloss.AdaptiveColor
// 	CanceledColor lipgloss.AdaptiveColor
// }

// Catppuccin color palette
var (
	// Base colors
	// rosewater = lipgloss.Color("#f5e0dc")
	// flamingo  = lipgloss.Color("#f2cdcd")
	// Pink  = lipgloss.Color("#f5c2e7")
	// Mauve = lipgloss.Color("#cba6f7")
	// Red   = lipgloss.Color("#f38ba8")
	// maroon    = lipgloss.Color("#eba0ac")
	// peach     = lipgloss.Color("#fab387")
	// Yellow = lipgloss.Color("#f9e2af")
	// Green  = lipgloss.Color("#a6e3a1")
	// Teal   = lipgloss.Color("#94e2d5")
	// sky       = lipgloss.Color("#89dceb")
	// sapphire  = lipgloss.Color("#74c7ec")
	// Blue     = lipgloss.Color("#89b4fa")
	// Lavender = lipgloss.Color("#b4befe")

	// Text colors
	// Text     = lipgloss.Color("#cdd6f4")
	// Subtext1 = lipgloss.Color("#bac2de")
	// subtext0 = lipgloss.Color("#a6adc8")
	// overlay2 = lipgloss.Color("#9399b2")
	// Overlay1 = lipgloss.Color("#7f849c")
	// overlay0 = lipgloss.Color("#6c7086")

	// Surface colors
	// surface2 = lipgloss.Color("#585b70")
	// surface1 = lipgloss.Color("#45475a")
	// surface0 = lipgloss.Color("#313244")
	// Base = lipgloss.Color("#1e1e2e")
	// mantle   = lipgloss.Color("#181825")
	// crust    = lipgloss.Color("#11111b")

	// ANSI colors standard and intense
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

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(MagentaANSI)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(LightWhiteANSI)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(RedANSI).
			Bold(true)

	// STATUS BAR
	StatusStyle = lipgloss.NewStyle().
			Foreground(BlackANSI).
			Background(LightGrayANSI).
			Padding(0, 1).
			Bold(false)

	StatusBarHeight = 1

	// WORKFLOWS
	ActiveWorkflowStyle = lipgloss.NewStyle().
				Foreground(GreenANSI)

	DisabledWorkflowStyle = lipgloss.NewStyle().
				Foreground(LightGrayANSI)

	// TABLE
	RowStyle = lipgloss.NewStyle().Bold(true).
			Foreground(WhiteANSI).Padding(0, 0, 0, 1)

	TableHeaderStyle = RowStyle.
				Bold(true)

	SelectedRowStyle = RowStyle.
				Background(YellowANSI)

	// JOBS
	SuccessStyle    = lipgloss.NewStyle().Foreground(GreenANSI)
	FailureStyle    = lipgloss.NewStyle().Foreground(RedANSI)
	CanceledStyle   = lipgloss.NewStyle().Foreground(YellowANSI)
	InProgressStyle = lipgloss.NewStyle().Foreground(CyanANSI)
)
