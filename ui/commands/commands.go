package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/pipeye/config"
	"github.com/cpaluszek/pipeye/github"
)

type ClientInitMsg struct {
	Client *github.Client
}

type ConfigInitMsg struct {
	Config *config.Config
}

type RepositoriesMsg struct {
	Repositories []*github.RepositoryData
}

type WorkflowsMsg struct {
	Workflows *github.RepositoryData
}

type WorkflowRunMsg struct {
	RunWithJobs *github.WorkflowRunWithJobs
}

type SectionChangedMsg struct{}

type ErrorMsg struct {
	Error error
}

// Commands
func InitConfig() tea.Msg {
	cfg, err := config.Load()
	if err != nil {
		return ErrorMsg{
			Error: err,
		}
	}
	return ConfigInitMsg{
		Config: cfg,
	}
}

func InitClient(token string) tea.Cmd {
	return func() tea.Msg {
		client, err := github.NewClient(token)
		if err != nil {
			return ErrorMsg{
				Error: err,
			}
		}
		return ClientInitMsg{
			Client: client,
		}
	}
}

func SectionChanged() tea.Msg {
	return SectionChangedMsg{}
}

func FetchRepositories(client *github.Client, names []string) tea.Cmd {
	return func() tea.Msg {
		repos, err := client.FetchRepositoriesWithWorkflows(names)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return RepositoriesMsg{
			Repositories: repos,
		}
	}
}

func GoToWorkflow(row github.RowData) tea.Cmd {
	return func() tea.Msg {
		workflows, ok := row.(*github.RepositoryData)
		if !ok {
			return ErrorMsg{
				Error: fmt.Errorf("selected row is not a repository"),
			}
		}
		return WorkflowsMsg{
			Workflows: workflows,
		}
	}
}

func GoToRun(row github.RowData) tea.Cmd {
	return func() tea.Msg {
		runWithJobs, ok := row.(*github.WorkflowRunWithJobs)
		if !ok {
			return ErrorMsg{
				Error: fmt.Errorf("selected row is not a workflow"),
			}
		}
		return WorkflowRunMsg{
			RunWithJobs: runWithJobs,
		}
	}
}

func OpenBrowser(url string) tea.Cmd {
	var cmd *exec.Cmd

	isWSL := false
	if runtime.GOOS == "linux" {
		_, hasWSLDistro := os.LookupEnv("WSL_DISTRO_NAME")
		_, hasWSLInterop := os.LookupEnv("WSL_INTEROP")
		isWSL = hasWSLDistro || hasWSLInterop
	}

	switch {
	case isWSL:
		cmd = exec.Command("explorer.exe", url)
	case runtime.GOOS == "darwin": // macOS
		cmd = exec.Command("open", url)
	case runtime.GOOS == "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return ErrorMsg{
				Error: fmt.Errorf("failed to open browser: %w", err),
			}
		}
		return nil
	})
}
