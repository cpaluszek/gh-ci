# gh-actions

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

1. Install the `gh` CLI - [instructions](https://github.com/cli/cli?tab=readme-ov-file#installation)
2. Install this extension:

```bash
gh extension install gh-actions
```

## Configuration


```yaml
github:
  repositories:
    - owner/repo1
    - owner/repo2
```

## Usage

```bash
gh actions
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
- [go-gh](https://github.com/cli/go-gh): Go library for the GitHub CLI

## Credits

This project was inspired by [gh-dash](https://github.com/dlvhdr/gh-dash) by [dlvhdr](https://github.com/dlvhdr) - a beautiful CLI dashboard for GitHub.
