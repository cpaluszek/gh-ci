package utils

import (
	"github.com/cpaluszek/gh-actions/ui/context"
)

func GetRunEventSymbol(ctx *context.Context, event string) string {
	switch event {
	case "pull_request":
		return ctx.Styles.PullRequest.Render(ctx.Theme.Symbols.PullRequest)
	case "push":
		return ctx.Styles.Push.Render(ctx.Theme.Symbols.Push)
	case "schedule":
		return ctx.Styles.Schedule.Render(ctx.Theme.Symbols.Schedule)
	case "release":
		return ctx.Styles.Tag.Render(ctx.Theme.Symbols.Tag)
	case "repository_dispatch":
		return ctx.Styles.WebHook.Render(ctx.Theme.Symbols.Webhook)
	case "workflow_dispatch", "dynamic":
		return ctx.Styles.Play.Render(ctx.Theme.Symbols.Play)
	case "fork":
		return ctx.Styles.Fork.Render(ctx.Theme.Symbols.Fork)
	case "deployment":
		return ctx.Styles.Deployment.Render(ctx.Theme.Symbols.Deployment)
	case "issue":
		return ctx.Styles.Issue.Render(ctx.Theme.Symbols.Issue)
	default:
		return ""
	}
}

func GetJobStatusSymbol(ctx *context.Context, status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return ctx.Styles.Success.Render(ctx.Theme.Symbols.JobSuccess)
		case "failure", "startup_failure", "timed_out", "action_required":
			return ctx.Styles.Failure.Render(ctx.Theme.Symbols.JobFailure)
		case "cancelled":
			return ctx.Styles.Canceled.Render(ctx.Theme.Symbols.JobCanceled)
		case "skipped":
			return ctx.Styles.Skipped.Render(ctx.Theme.Symbols.JobSkipped)
		case "neutral":
			return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
		default:
			return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
		}
	case "queued", "waiting", "pending", "requested":
		return ctx.Styles.InProgress.Render(ctx.Theme.Symbols.Queued)
	case "in_progress":
		return ctx.Styles.InProgress.Render(ctx.Theme.Symbols.InProgress)
	default:
		return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
	}
}

func GetStatusSymbol(ctx *context.Context, status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return ctx.Styles.Success.Render(ctx.Theme.Symbols.Success)
		case "failure", "timed_out", "startup_failure":
			return ctx.Styles.Failure.Render(ctx.Theme.Symbols.Failure)
		case "cancelled":
			return ctx.Styles.Canceled.Render(ctx.Theme.Symbols.Canceled)
		case "skipped":
			return ctx.Styles.Skipped.Render(ctx.Theme.Symbols.Skipped)
		case "neutral":
			return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
		default:
			return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
		}
	case "in_progress":
		return ctx.Styles.InProgress.Render(ctx.Theme.Symbols.InProgress)
	case "queued":
		return ctx.Styles.InProgress.Render(ctx.Theme.Symbols.Queued)
	default:
		return ctx.Styles.Default.Render(ctx.Theme.Symbols.Neutral)
	}
}
