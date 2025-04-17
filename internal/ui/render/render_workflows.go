package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/ui"
	gh "github.com/google/go-github/v71/github"
)

func RenderWorkflowsStatusBar(loading bool, style lipgloss.Style) string {
	var content string

	if loading {
		content = "Loading workflow... "
	} else {
		content = "↑/↓: navigate · esc/backspace: back to repositories"
	}

	return style.Render(content)
}

func RenderWorkflowsView(repo *gh.Repository, workflowsWithRuns []*github.WorkflowWithRuns, selectedRunIndex, width int, loading bool, err error) string {
	s := &strings.Builder{}

	// Header with repo info
	s.WriteString(ui.HeaderStyle.Render(fmt.Sprintf("\nRepository: %s\n", *repo.FullName)))

	// Handle errors
	if err != nil {
		s.WriteString(ui.ErrorTextStyle.Render(
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
		s.WriteString(ui.HeaderStyle.Render(
			fmt.Sprintf("\n%d. %s", i+1, workflow.GetName())))
		s.WriteString("\n   Path: " + workflow.GetPath())
		s.WriteString("\n   State: " + getWorkflowStateDisplay(workflow.GetState()))
		s.WriteString("\n   Created: " + formatTime(workflow.GetCreatedAt().Time))
		s.WriteString("     Last Updated: " + formatTime(workflow.GetUpdatedAt().Time))

		// Render runs if available
		if len(wwr.Runs) > 0 {
			s.WriteString("\n\n   Recent Runs:\n")
			s.WriteString(renderWorkflowRunsTable(wwr.Runs, selectedRunIndex, width))
		} else if wwr.Error != nil {
			s.WriteString("\n\n   " + ui.ErrorTextStyle.Render(
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
func renderWorkflowRunsTable(runsWithJobs []*github.WorkflowRunWithJobs, selectedRunIndex, width int) string {
	if len(runsWithJobs) == 0 {
		return "   No recent runs found."
	}

	s := &strings.Builder{}

	headers := []string{"Status", "Branch", "Triggered", "Duration", "Jobs", "Commit"}
	t := NewStyledTable(headers, width, selectedRunIndex)

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

		commitMsg := run.GetHeadCommit().GetMessage()

		// Style based on conclusion
		var statusStyle lipgloss.Style
		status := run.GetStatus()
		conclusion := run.GetConclusion()

		switch status {
		case "completed":
			// Use a nested switch for conclusion when status is completed
			switch conclusion {
			case "success":
				statusStyle = ui.SuccessStyle
			case "failure", "timed_out", "startup_failure":
				statusStyle = ui.FailureStyle
			case "cancelled", "skipped", "neutral":
				statusStyle = ui.CanceledStyle
			}
		case "in_progress":
			statusStyle = ui.InProgressStyle
			status = "running"
		}

		displayStatus := status
		if conclusion != "" && status == "completed" {
			displayStatus = conclusion
		}
		jobs := ""
		for _, job := range runWithJob.Jobs {
			switch job.GetConclusion() {
			case "success":
				statusStyle = ui.SuccessStyle
			case "failure":
				statusStyle = ui.FailureStyle
			case "cancelled":
				statusStyle = ui.CanceledStyle
			}

			jobs = lipgloss.JoinHorizontal(
				lipgloss.Center,
				jobs,
				statusStyle.Render("●"),
			)
		}

		// Table row
		var row = []string{
			statusStyle.Render(displayStatus),
			ui.RowStyle.Render(run.GetHeadBranch()),
			ui.RowStyle.Render(formatTime(run.GetCreatedAt().Time)),
			ui.RowStyle.Render(duration),
			ui.RowStyle.Render(jobs),
			ui.RowStyle.Render(commitMsg),
		}
		t.Row(row...)
	}
	s.WriteString(t.Render())

	return s.String()
}

func getWorkflowStateDisplay(state string) string {
	switch state {
	case "active":
		return ui.ActiveWorkflowStyle.Render("● active")
	case "disabled_manually", "disabled_inactivity":
		return ui.DisabledWorkflowStyle.Render("○ disabled")
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
