package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cpaluszek/pipeye/ui/commands"
	"github.com/cpaluszek/pipeye/ui/components/footer"
	"github.com/cpaluszek/pipeye/ui/context"
)

type Model struct {
	footer  footer.Model
	ctx     context.Context
	spinner spinner.Model
	loading bool
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	m := Model{
		footer:  footer.NewModel(),
		spinner: s,
		loading: true,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commands.InitConfig,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case commands.ConfigInitMsg:
		// TODO: handle error
		if msg.Error != nil {
			log.Println("Error loading config:", msg.Error)
			return m, tea.Quit
		}
		m.ctx.Config = msg.Config
		return m, commands.InitClient(m.ctx.Config.Github.Token)

	case commands.ClientInitMsg:
		if msg.Error != nil {
			log.Println("Error initializing client:", msg.Error)
			return m, tea.Quit
		}
		m.ctx.Client = msg.Client
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	}

	return m, nil
}

func (m Model) View() string {
	content := "Welcome to Pipeye!\n"
	if m.loading {
		return m.spinner.View() + "\n" + m.footer.View()
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Render(content),
		m.footer.View(),
	)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
}
