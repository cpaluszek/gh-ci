package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Colors  Colors
	Symbols Symbols
}

type Colors struct {
	Primary, Secondary, Faint                        lipgloss.AdaptiveColor
	PrimaryBorder, SecondaryBorder                   lipgloss.AdaptiveColor
	SelectedBackground, SelectedText                 lipgloss.AdaptiveColor
	Success, Warning, Error, Info, Progress, Skipped lipgloss.AdaptiveColor
	// Event colors
	PullRequest, Push, Schedule, Play, Issue lipgloss.AdaptiveColor
	Deployment, Tag, WebHook, Fork           lipgloss.AdaptiveColor
}

type Symbols struct {
	// Status symbols
	Success, Failure, Canceled, Skipped, Neutral, InProgress, Queued string
	// Job symbols
	JobSuccess, JobFailure, JobCanceled, JobSkipped, JobInProgress string
	// Event symbols
	PullRequest, Push, Schedule, Tag, Webhook, Fork, Deployment, Play, Issue string
}

var DefaultTheme = &Theme{
	Colors: Colors{
		Primary:            lipgloss.AdaptiveColor{Light: "000", Dark: "015"},
		Secondary:          lipgloss.AdaptiveColor{Light: "244", Dark: "251"},
		Faint:              lipgloss.AdaptiveColor{Light: "007", Dark: "254"},
		PrimaryBorder:      lipgloss.AdaptiveColor{Light: "013", Dark: "008"},
		SecondaryBorder:    lipgloss.AdaptiveColor{Light: "008", Dark: "007"},
		SelectedBackground: lipgloss.AdaptiveColor{Light: "006", Dark: "008"},
		SelectedText:       lipgloss.AdaptiveColor{Light: "000", Dark: "015"},
		Success:            lipgloss.AdaptiveColor{Light: "002", Dark: "002"},
		Warning:            lipgloss.AdaptiveColor{Light: "003", Dark: "003"},
		Error:              lipgloss.AdaptiveColor{Light: "001", Dark: "001"},
		Info:               lipgloss.AdaptiveColor{Light: "004", Dark: "004"},
		Progress:           lipgloss.AdaptiveColor{Light: "004", Dark: "004"},
		Skipped:            lipgloss.AdaptiveColor{Light: "007", Dark: "007"},
		PullRequest:        lipgloss.AdaptiveColor{Light: "004", Dark: "004"},
		Push:               lipgloss.AdaptiveColor{Light: "003", Dark: "003"},
		Schedule:           lipgloss.AdaptiveColor{Light: "007", Dark: "007"},
		Play:               lipgloss.AdaptiveColor{Light: "005", Dark: "005"},
		Issue:              lipgloss.AdaptiveColor{Light: "001", Dark: "001"},
		Deployment:         lipgloss.AdaptiveColor{Light: "002", Dark: "002"},
		Tag:                lipgloss.AdaptiveColor{Light: "006", Dark: "006"},
		WebHook:            lipgloss.AdaptiveColor{Light: "003", Dark: "003"},
		Fork:               lipgloss.AdaptiveColor{Light: "008", Dark: "008"},
	},
	Symbols: Symbols{
		Success:       "󰄬 ",
		Failure:       "󰅚 ",
		Canceled:      "󰔛 ",
		Skipped:       "󰒭 ",
		Neutral:       "󰘿 ",
		InProgress:    "󰑮 ",
		Queued:        "󰥔 ",
		JobSuccess:    "󰄯 ",
		JobFailure:    "󰅙 ",
		JobCanceled:   " ",
		JobSkipped:    " ",
		JobInProgress: "󱥸 ",
		PullRequest:   " ",
		Push:          " ",
		Schedule:      "󰃰 ",
		Tag:           " ",
		Webhook:       "󰛢 ",
		Fork:          " ",
		Deployment:    "󱓞 ",
		Play:          " ",
		Issue:         " ",
	},
}
