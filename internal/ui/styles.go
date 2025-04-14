package ui

import "github.com/charmbracelet/lipgloss"

// Catppuccin color palette
var (
	// Base colors
	rosewater = lipgloss.Color("#f5e0dc")
	flamingo  = lipgloss.Color("#f2cdcd")
	pink      = lipgloss.Color("#f5c2e7")
	mauve     = lipgloss.Color("#cba6f7")
	red       = lipgloss.Color("#f38ba8")
	maroon    = lipgloss.Color("#eba0ac")
	peach     = lipgloss.Color("#fab387")
	yellow    = lipgloss.Color("#f9e2af")
	green     = lipgloss.Color("#a6e3a1")
	teal      = lipgloss.Color("#94e2d5")
	sky       = lipgloss.Color("#89dceb")
	sapphire  = lipgloss.Color("#74c7ec")
	blue      = lipgloss.Color("#89b4fa")
	lavender  = lipgloss.Color("#b4befe")

	// Text colors
	text     = lipgloss.Color("#cdd6f4")
	subtext1 = lipgloss.Color("#bac2de")
	subtext0 = lipgloss.Color("#a6adc8")
	overlay2 = lipgloss.Color("#9399b2")
	overlay1 = lipgloss.Color("#7f849c")
	overlay0 = lipgloss.Color("#6c7086")

	// Surface colors
	surface2 = lipgloss.Color("#585b70")
	surface1 = lipgloss.Color("#45475a")
	surface0 = lipgloss.Color("#313244")
	base     = lipgloss.Color("#1e1e2e")
	mantle   = lipgloss.Color("#181825")
	crust    = lipgloss.Color("#11111b")

	// Application theme colors
	primaryColor   = blue
	secondaryColor = green
	accentColor    = mauve
	bgColor        = base
	textColor      = text
	errorColor     = red

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(pink)

	StatusStyle = lipgloss.NewStyle().
			Foreground(text).
			Background(base).
			Padding(0, 1).
			Bold(false)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lavender)

	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Left).
				Foreground(blue)

	RowStyle = lipgloss.NewStyle().
			Foreground(text)

	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(text).
				Background(base).
				Bold(true)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	ActiveWorkflowStyle = lipgloss.NewStyle().
				Foreground(green)

	DisabledWorkflowStyle = lipgloss.NewStyle().
				Foreground(overlay1)
)
