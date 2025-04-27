package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/cpaluszek/gh-ci/github"
	"github.com/cpaluszek/gh-ci/ui/context"
)

func FormatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

func formatDuration(d time.Duration) string {
	if d.Hours() >= 1 {
		return fmt.Sprintf("%.1fh", d.Hours())
	} else if d.Minutes() >= 1 {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
}

func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// CleanANSIEscapes removes ANSI reset sequences that cause rendering issues
// with lipgloss styled content when combined with other styles.
// This works around https://github.com/charmbracelet/lipgloss/issues/144
// where nested or adjacent styled content can have their styles reset
// by the automatic reset sequence (\x1b[0m) that lipgloss adds.
func CleanANSIEscapes(s string) string {
	return strings.ReplaceAll(s, "\x1b[0m", "")
}

func GetWorkflowRunDuration(wr *github.WorkflowRun) string {
	if wr == nil {
		return ""
	}
	var duration string
	if wr.UpdatedAt.After(wr.CreatedAt) {
		durationTime := wr.UpdatedAt.Sub(wr.CreatedAt)
		duration = formatDuration(durationTime)
	} else {
		duration = "running"
	}
	return duration
}

func GetWorkflowRunStatus(ctx *context.Context, wr *github.WorkflowRun) string {
	if wr == nil {
		return ""
	}
	status := wr.Status
	conclusion := wr.Conclusion
	statusSymbol := GetStatusSymbol(ctx, status, conclusion)
	content := ""
	if conclusion != "" && status == "completed" {
		content = statusSymbol + conclusion
	} else if status == "in_progress" {
		content = statusSymbol + "running"
	} else {
		content = statusSymbol + status
	}

	return CleanANSIEscapes(content)
}

func GetJobDuration(job *github.Job) string {
	if job == nil {
		return ""
	}
	var duration string
	if job.CompletedAt.After(job.StartedAt) {
		durationTime := job.CompletedAt.Sub(job.StartedAt)
		duration = formatDuration(durationTime)
	} else {
		duration = "running"
	}
	return duration
}
