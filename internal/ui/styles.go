package ui

import "github.com/charmbracelet/lipgloss"

var (
    StatusStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Background(lipgloss.Color("#333333")).
        Padding(0, 1).
        Bold(false)

	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#76ABDF"))

	TableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Left)

	RowStyle = lipgloss.NewStyle()

	SelectedRowStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("15")).
        Background(lipgloss.Color("57")).
        Bold(true)
)

