package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/config"
	"github.com/cpaluszek/gh-ci/ui"
)

func Execute() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Redirect logs to a file
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		panic(err)
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	model := ui.NewModel(cfg)

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to run program:", err)
	}
}
