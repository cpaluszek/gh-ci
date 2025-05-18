package styles

import "github.com/charmbracelet/lipgloss"
import bbhelp "github.com/charmbracelet/bubbles/help"

type Styles struct {
	Header           lipgloss.Style
	Title            lipgloss.Style
	Error            lipgloss.Style
	Footer        lipgloss.Style
    Help            bbhelp.Styles
	Spinner          lipgloss.Style
	Row              lipgloss.Style
	TableHeader      lipgloss.Style
	SelectedRow      lipgloss.Style
	SectionContainer lipgloss.Style
	SideBar          lipgloss.Style
	Success          lipgloss.Style
	Failure          lipgloss.Style
	Canceled         lipgloss.Style
	InProgress       lipgloss.Style
	Skipped          lipgloss.Style
	Default          lipgloss.Style
	PullRequest      lipgloss.Style
	Push             lipgloss.Style
	Schedule         lipgloss.Style
	Play             lipgloss.Style
	Issue            lipgloss.Style
	Deployment       lipgloss.Style
	Tag              lipgloss.Style
	WebHook          lipgloss.Style
	Fork             lipgloss.Style
}

func BuildStyles(theme Theme) Styles {
	var s Styles

	s.Header = lipgloss.NewStyle().Bold(true).
		Border(lipgloss.NormalBorder()).BorderForeground(theme.Colors.PrimaryBorder).
		BorderBottom(true).BorderTop(false).BorderLeft(false).BorderRight(false)

	s.Title = lipgloss.NewStyle().
		Foreground(theme.Colors.Primary).
		Bold(true)

	s.Error = lipgloss.NewStyle().
		Foreground(theme.Colors.Error).
		Bold(true)

    helpText := lipgloss.NewStyle()
    helpKeyText := lipgloss.NewStyle()
	s.Help = bbhelp.Styles{
		ShortDesc:      helpText.Foreground(theme.Colors.Faint),
		FullDesc:       helpText.Foreground(theme.Colors.Faint),
		ShortSeparator: helpText.Foreground(theme.Colors.SecondaryBorder),
		FullSeparator:  helpText.Foreground(theme.Colors.SecondaryBorder),
		FullKey:        helpKeyText,
		ShortKey:       helpKeyText,
		Ellipsis:       helpText,
	}
	s.Footer = lipgloss.NewStyle().Padding(0, 1).
		Border(lipgloss.NormalBorder()).BorderForeground(theme.Colors.PrimaryBorder).
		BorderBottom(false).BorderTop(true).BorderLeft(false).BorderRight(false)

	s.Spinner = lipgloss.NewStyle().Bold(true).Padding(0, 1)

	s.Row = lipgloss.NewStyle()

	s.TableHeader = lipgloss.NewStyle().Bold(true)

	s.SelectedRow = lipgloss.NewStyle().Background(theme.Colors.SelectedBackground).Foreground(theme.Colors.SelectedText)

	s.SectionContainer = lipgloss.NewStyle().Padding(0, 1)

	s.SideBar = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(theme.Colors.PrimaryBorder).
		BorderBottom(false).
		BorderTop(false).
		BorderLeft(true).
		BorderRight(false)

	s.Success = lipgloss.NewStyle().Foreground(theme.Colors.Success)

	s.Failure = lipgloss.NewStyle().Foreground(theme.Colors.Error)

	s.Canceled = lipgloss.NewStyle().Foreground(theme.Colors.Warning)

	s.InProgress = lipgloss.NewStyle().Foreground(theme.Colors.Progress)

	s.Skipped = lipgloss.NewStyle().Foreground(theme.Colors.Skipped)

	s.Default = lipgloss.NewStyle().Foreground(theme.Colors.Primary)

	s.PullRequest = lipgloss.NewStyle().Foreground(theme.Colors.PullRequest).Bold(true)

	s.Push = lipgloss.NewStyle().Foreground(theme.Colors.Push).Bold(true)

	s.Schedule = lipgloss.NewStyle().Foreground(theme.Colors.Schedule).Bold(true)

	s.Play = lipgloss.NewStyle().Foreground(theme.Colors.Play).Bold(true)

	s.Issue = lipgloss.NewStyle().Foreground(theme.Colors.Issue).Bold(true)

	s.Deployment = lipgloss.NewStyle().Foreground(theme.Colors.Deployment).Bold(true)

	s.Tag = lipgloss.NewStyle().Foreground(theme.Colors.Tag).Bold(true)

	s.WebHook = lipgloss.NewStyle().Foreground(theme.Colors.WebHook).Bold(true)

	s.Fork = lipgloss.NewStyle().Foreground(theme.Colors.Fork).Bold(true)

	return s
}
