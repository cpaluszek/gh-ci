package cmd

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/ui"
)

func Execute() {
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

	model := ui.NewModel()

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal("Failed to run program:", err)
	}
}
