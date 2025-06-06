package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cpaluszek/gh-ci/config"
	"github.com/cpaluszek/gh-ci/github"
)

type ClientInitMsg struct {
	Client *github.Client
}

type ConfigInitMsg struct {
	Config *config.Config
}

type RepositoriesMsg struct {
	Repositories []*github.Repository
}

type WorkflowsMsg struct {
	Workflows *github.Repository
}

type WorkflowRunMsg struct {
	RunWithJobs *github.WorkflowRun
}

type LogsMsg struct {
	Steps []github.Steplog
}

type GotostepMsg struct {
	RunWithJobs *github.Job
}

type SectionChangedMsg struct{}

type ErrorMsg struct {
	Error error
}

// Commands
func InitClient() tea.Cmd {
	return func() tea.Msg {
		client, err := github.NewClient()
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

func FetchStepLogs(client *github.Client, job *github.Job) tea.Cmd {
	return func() tea.Msg {
		if job == nil {
			return ErrorMsg{
				Error: fmt.Errorf("workflow run is nil"),
			}
		}
		info, err := github.ParseGitHubURL(job.GetURL())
		if err != nil {
			return ErrorMsg{Error: err}
		}
		steps, err := client.GetLogs(info.User, info.Repo, info.RunID, "1", job.Name)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return LogsMsg{
			Steps: steps,
		}
	}
}

func FetchLogs(client *github.Client, job *github.Job) tea.Cmd {
	return func() tea.Msg {
		if job == nil {
			return ErrorMsg{
				Error: fmt.Errorf("workflow run is nil"),
			}
		}
		info, err := github.ParseGitHubURL(job.GetURL())
		if err != nil {
			return ErrorMsg{Error: err}
		}
		steps, err := client.GetLogs(info.User, info.Repo, info.RunID, "1", job.Name)
		if err != nil {
			return ErrorMsg{Error: err}
		}
		return LogsMsg{
			Steps: steps,
		}
	}
}

func GoToStep(row github.RowData) tea.Cmd {
	return func() tea.Msg {
		runWithJobs, ok := row.(*github.Job)
		if !ok {
			return ErrorMsg{
				Error: fmt.Errorf("selected row is not a workflow"),
			}
		}
		return GotostepMsg{
			RunWithJobs: runWithJobs,
		}
	}
}

func GoToWorkflow(row github.RowData) tea.Cmd {
	return func() tea.Msg {
		workflows, ok := row.(*github.Repository)
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
		runWithJobs, ok := row.(*github.WorkflowRun)
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
