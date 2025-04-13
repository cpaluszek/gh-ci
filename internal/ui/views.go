package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v71/github"
)

func RenderRepositoriesTable(repositories []*github.Repository, width int) string {
	var s strings.Builder

	nameWidth, langWidth, starsWidth, updatedWidth, workflowsWidth := calculateColumnWidths(width)
	totalWidth := nameWidth + langWidth + starsWidth + updatedWidth + workflowsWidth

	s.WriteString(HeaderStyle.Render("GitHub Repositories with Workflows"))
	s.WriteString("\n\n")

	// Column headers
	headers := lipgloss.JoinHorizontal(lipgloss.Top,
		TableHeaderStyle.Width(nameWidth).Align(lipgloss.Left).Render("Repository"),
		TableHeaderStyle.Width(langWidth).Align(lipgloss.Left).Render("Language"),
		TableHeaderStyle.Width(starsWidth).Align(lipgloss.Left).Render("Stars"),
		TableHeaderStyle.Width(updatedWidth).Align(lipgloss.Left).Render("Last Updated"),
		TableHeaderStyle.Width(workflowsWidth).Align(lipgloss.Left).Render("Workflows"),
		)
	s.WriteString(headers + "\n")
	s.WriteString(strings.Repeat("─", totalWidth) + "\n")

	for _, repo := range repositories {
		language := ""
		if repo.Language != nil {
			language = *repo.Language
		}

		stars := "0"
		if repo.StargazersCount != nil {
			stars = fmt.Sprintf("%d", *repo.StargazersCount)
		}

		updated := "Unknown"
		if repo.UpdatedAt != nil {
			updated = repo.UpdatedAt.Format("Jan 2, 2006")
		}

		// Create row style with alternating background
		rowStyle := lipgloss.NewStyle()
		// if i % 2 == 1 {
		// 	rowStyle = rowStyle.Background(lipgloss.Color("0"))
		// }

		row := lipgloss.JoinHorizontal(lipgloss.Top,
			rowStyle.Width(nameWidth).Align(lipgloss.Left).Render(*repo.FullName),
			rowStyle.Width(langWidth).Align(lipgloss.Left).Render(language),
			rowStyle.Width(starsWidth).Align(lipgloss.Left).Render(stars),
			rowStyle.Width(updatedWidth).Align(lipgloss.Left).Render(updated),
			rowStyle.Width(workflowsWidth).Align(lipgloss.Left).Render("✓"),
			)
		s.WriteString(row + "\n")
	}

	return s.String()
}

func RenderStatusBar(loading bool, repoCount int, style lipgloss.Style) string {
	var content string

	if loading {
		content = "Loading repositories... "
	} else if repoCount > 0 {
		content = fmt.Sprintf("Found %d repositories with workflows · q: quit", repoCount)
	} else {
		content = "No repositories found · q: quit"
	}

	return style.Render(content)
}

func calculateColumnWidths(width int) (nameWidth, langWidth, starsWidth, updatedWidth, workflowsWidth int) {
	availableWidth := width
	if availableWidth == 0 {
		availableWidth = 100 // Fallback
	}

	nameWidth = int(float64(availableWidth) * 0.4)
	langWidth = int(float64(availableWidth) * 0.15)
	starsWidth = int(float64(availableWidth) * 0.10)
	updatedWidth = int(float64(availableWidth) * 0.25)
	workflowsWidth = int(float64(availableWidth) * 0.10)

	// Ensure minimum widths
	nameWidth = max(nameWidth, 20)
	langWidth = max(langWidth, 10)
	starsWidth = max(starsWidth, 6)
	updatedWidth = max(updatedWidth, 15)
	workflowsWidth = max(workflowsWidth, 8)

	// Adjust if total exceeds available width
	totalWidth := nameWidth + langWidth + starsWidth + updatedWidth + workflowsWidth
	if totalWidth > availableWidth {
		ratio := float64(availableWidth) / float64(totalWidth)
		nameWidth = int(float64(nameWidth) * ratio)
		langWidth = int(float64(langWidth) * ratio)
		starsWidth = int(float64(starsWidth) * ratio)
		updatedWidth = int(float64(updatedWidth) * ratio)
		workflowsWidth = int(float64(workflowsWidth) * ratio)
	}

	return
}
