package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cpaluszek/pipeye/internal/config"
	"github.com/cpaluszek/pipeye/internal/github"
	"github.com/cpaluszek/pipeye/internal/models"
	"github.com/cpaluszek/pipeye/internal/ui/views"
)

type Model struct {
	views.BaseView
	config     *config.Config
	Client     *github.Client
	ListView   views.ListView
	detailView views.DetailView
	showDetail bool
}

func New(cfg *config.Config) *Model {
	return &Model{
		config:     cfg,
		showDetail: false,
	}
}

func (m Model) Init() tea.Cmd {
	return models.InitClient(m.config.Github.Token)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.showDetail {
			detailView, cmd := m.detailView.Update(msg)
			m.detailView = detailView
			cmds = append(cmds, cmd)
		} else {
			listView, cmd := m.ListView.Update(msg)
			m.ListView = listView
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		if m.showDetail {
			// Handle key events in detail view
			switch msg.String() {
			case "esc", "backspace":
				m.showDetail = false
				return m, nil
			default:
				var cmd tea.Cmd
				m.detailView, cmd = m.detailView.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			// Handle key events in list view
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				if m.ListView.HasRepositories() {
					selectedRepo := m.ListView.GetSelectedRepo()
					m.detailView = views.NewDetailView(selectedRepo, m.ListView.Viewport, m.Client)
					m.showDetail = true
					return m, m.detailView.Init()
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.ListView, cmd = m.ListView.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case models.ClientInitializedMsg:
		m.Client = msg.Client
		m.ListView = views.NewListView(m.Client)
		return m, m.ListView.Init()

	case models.RepositoriesMsg:
		var cmd tea.Cmd
		m.ListView, cmd = m.ListView.Update(msg)
		cmds = append(cmds, cmd)

	case models.DetailViewMsg:
		if m.showDetail {
			var cmd tea.Cmd
			m.detailView, cmd = m.detailView.Update(msg)
			cmds = append(cmds, cmd)
		}

	case models.ErrMsg:
		if m.showDetail {
			m.detailView.Error = msg.Err
			m.detailView.Loading = false
		} else {
			m.ListView.Error = msg.Err
			m.ListView.Loading = false
		}
		return m, nil

	default:
		if !m.showDetail {
			var cmd tea.Cmd
			m.ListView, cmd = m.ListView.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.detailView, cmd = m.detailView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() string {
	if m.showDetail {
		return m.detailView.View()
	}
	return m.ListView.View()
}

func (m *Model) Run() error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	_, error := p.Run()
	return error
}
