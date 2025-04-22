package utils

import (
	"fmt"
	"strings"
	"time"
)

func FormatTime(t time.Time) string {
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

// CleanANSIEscapes removes ANSI reset sequences that cause rendering issues
// with lipgloss styled content when combined with other styles.
// This works around https://github.com/charmbracelet/lipgloss/issues/144
// where nested or adjacent styled content can have their styles reset
// by the automatic reset sequence (\x1b[0m) that lipgloss adds.
func CleanANSIEscapes(s string) string {
	return strings.ReplaceAll(s, "\x1b[0m", "")
}
