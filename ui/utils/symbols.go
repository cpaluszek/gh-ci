package utils

import "github.com/cpaluszek/pipeye/ui/styles"

func GetRunEventSymbol(event string) string {
	switch event {
	case "pull_request":
		return styles.PullRequestStyle.Render(styles.PullRequestSymbol)
	case "push":
		return styles.PushStyle.Render(styles.PushSymbol)
	case "schedule":
		return styles.ScheduleStyle.Render(styles.ScheduleSymbol)
	case "release":
		return styles.TagStyle.Render(styles.TagStmbol)
	case "repository_dispatch":
		return styles.WebHookStyle.Render(styles.WebhookSymbol)
	case "workflow_dispatch", "dynamic":
		return styles.PlayStyle.Render(styles.PlaySymbol)
	case "fork":
		return styles.ForkStyle.Render(styles.ForkSymbol)
	case "deployment":
		return styles.DeploymentStyle.Render(styles.DeploymentSymbol)
	case "issue":
		return styles.IssueStyle.Render(styles.IssueSymbol)
	default:
		return ""
	}
}

func GetJobStatusSymbol(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return styles.SuccessStyle.Render(styles.JobSuccessDot)
		case "failure", "startup_failure", "timed_out", "action_required":
			return styles.FailureStyle.Render(styles.JobFailureDot)
		case "cancelled":
			return styles.CanceledStyle.Render(styles.JobCanceledDot)
		case "skipped":
			return styles.SkippedStyle.Render(styles.JobSkippedDot)
		case "neutral":
			return styles.DefaultStyle.Render(styles.NeutralSymbol)
		default:
			return styles.DefaultStyle.Render(styles.NeutralSymbol)
		}
	case "queued", "waiting", "pending", "requested":
		return styles.InProgressStyle.Render(styles.QueuedSymbol)
	case "in_progress":
		return styles.InProgressStyle.Render(styles.InProgressSymbol)
	default:
		return styles.DefaultStyle.Render(styles.NeutralSymbol)
	}
}

func GetStatusSymbol(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return styles.SuccessStyle.Render(styles.SuccessSymbol)
		case "failure", "timed_out", "startup_failure":
			return styles.FailureStyle.Render(styles.FailureSymbol)
		case "cancelled":
			return styles.CanceledStyle.Render(styles.CanceledSymbol)
		case "skipped":
			return styles.SkippedStyle.Render(styles.SkippedSymbol)
		case "neutral":
			return styles.DefaultStyle.Render(styles.NeutralSymbol)
		default:
			return styles.DefaultStyle.Render(styles.NeutralSymbol)
		}
	case "in_progress":
		return styles.InProgressStyle.Render(styles.InProgressSymbol)
	case "queued":
		return styles.InProgressStyle.Render(styles.QueuedSymbol)
	default:
		return styles.DefaultStyle.Render(styles.NeutralSymbol)
	}
}
