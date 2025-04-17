package ui

import "github.com/charmbracelet/lipgloss"

// Catppuccin color palette
var (
	// Base colors
	// rosewater = lipgloss.Color("#f5e0dc")
	// flamingo  = lipgloss.Color("#f2cdcd")
	Pink  = lipgloss.Color("#f5c2e7")
	Mauve = lipgloss.Color("#cba6f7")
	Red   = lipgloss.Color("#f38ba8")
	// maroon    = lipgloss.Color("#eba0ac")
	// peach     = lipgloss.Color("#fab387")
	Yellow = lipgloss.Color("#f9e2af")
	Green  = lipgloss.Color("#a6e3a1")
	Teal   = lipgloss.Color("#94e2d5")
	// sky       = lipgloss.Color("#89dceb")
	// sapphire  = lipgloss.Color("#74c7ec")
	Blue     = lipgloss.Color("#89b4fa")
	Lavender = lipgloss.Color("#b4befe")

	// Text colors
	Text     = lipgloss.Color("#cdd6f4")
	Subtext1 = lipgloss.Color("#bac2de")
	// subtext0 = lipgloss.Color("#a6adc8")
	// overlay2 = lipgloss.Color("#9399b2")
	Overlay1 = lipgloss.Color("#7f849c")
	// overlay0 = lipgloss.Color("#6c7086")

	// Surface colors
	// surface2 = lipgloss.Color("#585b70")
	// surface1 = lipgloss.Color("#45475a")
	// surface0 = lipgloss.Color("#313244")
	Base = lipgloss.Color("#1e1e2e")
	// mantle   = lipgloss.Color("#181825")
	// crust    = lipgloss.Color("#11111b")

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Pink)

	StatusStyle = lipgloss.NewStyle().
			Foreground(Text).
			Background(Base).
			Padding(0, 1).
			Bold(false)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	ActiveWorkflowStyle = lipgloss.NewStyle().
				Foreground(Green)

	DisabledWorkflowStyle = lipgloss.NewStyle().
				Foreground(Overlay1)

	// TABLE
	RowStyle = lipgloss.NewStyle().
			Foreground(Text).Padding(0, 0, 0, 1)

	TableHeaderStyle = RowStyle.
				Bold(true)

	SelectedRowStyle = RowStyle.
				Background(Base)
)
