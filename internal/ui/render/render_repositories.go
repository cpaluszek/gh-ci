package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/internal/ui"
	gh "github.com/google/go-github/v71/github"
)

func RenderRepositoriesStatusBar(loading bool, repoCount int, style lipgloss.Style) string {
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

func RenderRepositoriesTable(repositories []*gh.Repository, selectedIndex int, width int) string {
	var s strings.Builder

	s.WriteString("\n\n")

	headers := []string{"Repository", "Language", "Stars", "Last Updated"}
	t := NewStyledTable(headers, width, selectedIndex)

	for _, repo := range repositories {
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

		var row = []string{
			ui.RowStyle.Render(*repo.FullName),
			ui.RowStyle.Render(language),
			ui.RowStyle.Render(stars),
			ui.RowStyle.Render(updated),
		}
		t.Row(row...)
	}

	s.WriteString(t.Render())

	return s.String()
}
