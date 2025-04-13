package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v71/github"
)

func RenderRepositoriesTable(repositories []*github.Repository, selectedIndex int, width int) string {
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
            rowStyle = SelectedRowStyle
        } else {
            rowStyle = RowStyle
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

func RenderDetailViewStatusBar(loading bool, style lipgloss.Style) string {
	var content string

	if loading {
		content = "Loading workflow... "
	} else {
		content = "Workflow · <esc>: close"
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

func RenderDetailView(repo *github.Repository, workflows []*github.Workflow, loading bool, err error) string {
	var sb strings.Builder

	// Repository header
	repoName := *repo.FullName
	repoHeader := HeaderStyle.Render(repoName)
	sb.WriteString(repoHeader + "\n\n")

	// Basic info
	sb.WriteString(fmt.Sprintf("Description: %s\n", stringOrEmpty(repo.Description)))
	sb.WriteString(fmt.Sprintf("Language: %s\n", stringOrEmpty(repo.Language)))
	sb.WriteString(fmt.Sprintf("Stars: %d\n", intOrZero(repo.StargazersCount)))
	sb.WriteString(fmt.Sprintf("Forks: %d\n", intOrZero(repo.ForksCount)))
	sb.WriteString(fmt.Sprintf("Watchers: %d\n", intOrZero(repo.WatchersCount)))
	sb.WriteString(fmt.Sprintf("Open Issues: %d\n", intOrZero(repo.OpenIssuesCount)))
	if repo.UpdatedAt != nil && !repo.UpdatedAt.IsZero() {
		sb.WriteString(fmt.Sprintf("Last Updated: %s\n", formatTime(repo.UpdatedAt.Time)))
	}
	sb.WriteString("\n")

	// Workflows section
	sb.WriteString(HeaderStyle.Render("Workflows") + "\n\n")

	if loading {
		sb.WriteString("Loading workflows...\n")
	} else if err != nil {
		sb.WriteString(fmt.Sprintf("Error loading workflows: %v\n", err))
	} else if len(workflows) == 0 {
		sb.WriteString("No workflows found for this repository.\n")
	} else {
		for _, wf := range workflows {
			sb.WriteString(fmt.Sprintf("• %s\n", *wf.Name))
			sb.WriteString(fmt.Sprintf("  Path: %s\n", *wf.Path))
			sb.WriteString(fmt.Sprintf("  State: %s\n", *wf.State))
			sb.WriteString(fmt.Sprintf("  BadgeUrl: %s\n", *wf.BadgeURL))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\nPress ESC to return to repository list")

	return sb.String()
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return "None"
	}
	return *s
}

func intOrZero(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, pluralize(days))
	default:
		return t.Format("Jan 2, 2006")
	}
}

func pluralize(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
