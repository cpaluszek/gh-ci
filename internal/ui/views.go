package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	github "github.com/cpaluszek/pipeye/internal/github_client"
	gh "github.com/google/go-github/v71/github"
)

func RenderRepositoriesTable(repositories []*gh.Repository, selectedIndex int, width int) string {
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
		content = "↑/↓: navigate · f/b: page up/down · esc: back to repositories"
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

func RenderDetailView(repo *gh.Repository, workflowsWithRuns []*github.WorkflowWithRuns, loading bool, err error) string {
	s := &strings.Builder{}

	// Header with repo info
	s.WriteString(HeaderStyle.Render(fmt.Sprintf("Repository: %s\n\n", *repo.FullName)))

	// Handle errors
	if err != nil {
		s.WriteString(ErrorTextStyle.Render(
			fmt.Sprintf("Error loading workflows: %v\n", err)))
		return s.String()
	}

	if loading {
		s.WriteString("Loading workflows...\n")
		return s.String()
	}

	if len(workflowsWithRuns) == 0 {
		s.WriteString("No workflows found for this repository.\n")
		return s.String()
	}

	// Render each workflow with its runs
	for i, wwr := range workflowsWithRuns {
		workflow := wwr.Workflow

		// Workflow header
		s.WriteString(HeaderStyle.Render(
			fmt.Sprintf("\n%d. %s", i+1, workflow.GetName())))
		s.WriteString("\n   Path: " + workflow.GetPath())
		s.WriteString("\n   State: " + getWorkflowStateDisplay(workflow.GetState()))
		s.WriteString("\n   Created: " + formatTime(workflow.GetCreatedAt().Time))
		s.WriteString("\n   Last Updated: " + formatTime(workflow.GetUpdatedAt().Time))

		// Render runs if available
		if len(wwr.Runs) > 0 {
			s.WriteString("\n\n   Recent Runs:\n")
			s.WriteString(renderWorkflowRunsTable(wwr.Runs))
		} else if wwr.Error != nil {
			s.WriteString("\n\n   " + ErrorTextStyle.Render(
				fmt.Sprintf("Error loading runs: %v", wwr.Error)))
		} else {
			s.WriteString("\n\n   No recent runs found.")
		}

		// Add separator between workflows
		if i < len(workflowsWithRuns)-1 {
			s.WriteString("\n\n" + strings.Repeat("─", 50) + "\n")
		}
	}

	return s.String()
}

// renderWorkflowRunsTable renders a table of workflow runs
func renderWorkflowRunsTable(runsWithJobs []*github.WorkflowRunWithJobs) string {
	if len(runsWithJobs) == 0 {
		return "   No recent runs found."
	}

	s := &strings.Builder{}

	// Table header
	s.WriteString("   " + lipgloss.JoinHorizontal(lipgloss.Top,
		TableHeaderStyle.Width(15).Align(lipgloss.Left).Render("Status"),
		TableHeaderStyle.Width(10).Align(lipgloss.Left).Render("Branch"),
		TableHeaderStyle.Width(20).Align(lipgloss.Left).Render("Triggered"),
		TableHeaderStyle.Width(15).Align(lipgloss.Left).Render("Duration"),
	) + "\n")

	// Table rows
	for _, runWithJob := range runsWithJobs {
		run := runWithJob.Run
		// Calculate duration
		var duration string
		if run.GetUpdatedAt().After(run.GetCreatedAt().Time) {
			durationTime := run.GetUpdatedAt().Sub(run.GetCreatedAt().Time)
			if durationTime.Hours() >= 1 {
				duration = fmt.Sprintf("%.1fh", durationTime.Hours())
			} else if durationTime.Minutes() >= 1 {
				duration = fmt.Sprintf("%.1fm", durationTime.Minutes())
			} else {
				duration = fmt.Sprintf("%.1fs", durationTime.Seconds())
			}
		} else {
			duration = "running"
		}

		// Style based on conclusion
		statusStyle := RowStyle.Width(15)
		status := run.GetStatus()
		conclusion := run.GetConclusion()

		switch status {
		case "completed":
			// Use a nested switch for conclusion when status is completed
			switch conclusion {
			case "success":
				statusStyle = statusStyle.Foreground(lipgloss.Color("10")) // Green
			case "failure", "timed_out":
				statusStyle = statusStyle.Foreground(lipgloss.Color("9")) // Red
			case "cancelled", "skipped", "neutral":
				statusStyle = statusStyle.Foreground(lipgloss.Color("11")) // Yellow
			}
		case "in_progress":
			statusStyle = statusStyle.Foreground(lipgloss.Color("14")) // Cyan
			status = "running"
		}

		displayStatus := status
		if conclusion != "" && status == "completed" {
			displayStatus = conclusion
		}
		jobHeader := ""
		for _, job := range runWithJob.Jobs {
			statusStyle := RowStyle
			switch job.GetConclusion() {
			case "success":
				statusStyle = statusStyle.Foreground(lipgloss.Color("10")) // Green
			case "failure":
				statusStyle = statusStyle.Foreground(lipgloss.Color("9")) // Red
			case "cancelled":
				statusStyle = statusStyle.Foreground(lipgloss.Color("11")) // Yellow
			}

			jobHeader = lipgloss.JoinHorizontal(
				lipgloss.Center,
				statusStyle.Render("●"),
				" ",
			)
		}

		// Table row
		s.WriteString("   " + lipgloss.JoinHorizontal(lipgloss.Top,
			statusStyle.Render(displayStatus),
			RowStyle.Width(10).Render(run.GetHeadBranch()),
			RowStyle.Width(20).Render(formatTime(run.GetCreatedAt().Time)),
			RowStyle.Width(15).Render(duration),
			RowStyle.Render(jobHeader),
		) + "\n")
	}

	return s.String()
}

func getWorkflowStateDisplay(state string) string {
	switch state {
	case "active":
		return ActiveWorkflowStyle.Render("● active")
	case "disabled_manually", "disabled_inactivity":
		return DisabledWorkflowStyle.Render("○ disabled")
	default:
		return state
	}
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
