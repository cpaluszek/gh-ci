package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/ui"
	gh "github.com/google/go-github/v71/github"
)

func RenderRepositoriesTable(repositories []*gh.Repository, selectedIndex int, width int) string {
	var s strings.Builder

	nameWidth, langWidth, starsWidth, updatedWidth, workflowsWidth := calculateColumnWidths(width)
	totalWidth := nameWidth + langWidth + starsWidth + updatedWidth + workflowsWidth

	s.WriteString(ui.HeaderStyle.Render("GitHub Repositories with Workflows"))
	s.WriteString("\n\n")

	// Column headers
	headers := lipgloss.JoinHorizontal(lipgloss.Top,
		ui.TableHeaderStyle.Width(nameWidth).Align(lipgloss.Left).Render("Repository"),
		ui.TableHeaderStyle.Width(langWidth).Align(lipgloss.Left).Render("Language"),
		ui.TableHeaderStyle.Width(starsWidth).Align(lipgloss.Left).Render("Stars"),
		ui.TableHeaderStyle.Width(updatedWidth).Align(lipgloss.Left).Render("Last Updated"),
		ui.TableHeaderStyle.Width(workflowsWidth).Align(lipgloss.Left).Render("Workflows"),
	)
	s.WriteString(headers + "\n")
	s.WriteString(strings.Repeat("─", totalWidth) + "\n")

	for i, repo := range repositories {
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

		var rowStyle lipgloss.Style
		if i == selectedIndex {
			rowStyle = ui.SelectedRowStyle
		} else {
			rowStyle = ui.RowStyle
		}

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
		content = fmt.Sprintf("%d repositories · ↑/↓: navigate · enter: select · q: quit", repoCount)
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
