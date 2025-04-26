# Pipeye

A terminal UI application for monitoring your GitHub Actions workflows and runs.

## Features

- 🔍 Browse your GitHub repositories with Actions workflows
- 📊 View workflows and their recent runs
- 🔄 Monitor run status in real-time with visual indicators (WIP)
- 👁️ See job status and details for each workflow run (WIP)

## Requirements

### Fonts
Pipeye uses Nerd Font icons for workflow and job status indicators. For the best experience:

1. Install a [Nerd Font](https://www.nerdfonts.com/font-downloads) compatible font of your choice
2. Configure your terminal to use the installed Nerd Font

Without a Nerd Font, the status icons will appear as placeholder characters or missing glyphs.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/cpaluszek/pipeye.git
cd pipeye

# Build the application
go build -o pipeye .

# Move to a directory in your PATH (optional)
sudo mv pipeye /usr/local/bin/
```

## Configuration

Pipeye requires a GitHub personal access token to access your repositories and workflows. Create a config.yaml file in the same directory where you run the application:

```yaml
github:
  token: your-github-token
  repositories:
    - owner/repo1
    - owner/repo2
```

## Usage

```bash
# Start pipeye
pipeye
```

### Navigation

- **Repository List:**
  - `↑`/`k`, `↓`/`j`: Navigate through repositories
  - `Enter`: Select repository and view its workflows
  - `o`: Open selected repository in browser
  - `q`: Quit

- **Workflow View:**
  - `↑`/`k`, `↓`/`j`: Navigate through workflow runs
  - `Enter`: Select run and view jobs 
  - `o`: Open selected workflow run in browser
  - `Esc`/`Backspace`: Return to repository list

- **Workflow View:**
  - `↑`/`k`, `↓`/`j`: Navigate through workflow runs
  - `o`: Open selected workflow run in browser
  - `Esc`/`Backspace`: Return to workflow list

## Development

Pipeye is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions for terminal applications
- [go-github](https://github.com/google/go-github): GitHub API client for Go

## Credits

This project was inspired by [gh-dash](https://github.com/dlvhdr/gh-dash) by [dlvhdr](https://github.com/dlvhdr) - a beautiful CLI dashboard for GitHub.
