# gh-ci

A terminal UI application for monitoring your GitHub Actions workflows and runs.

![Demo GIF](./assets/demo.gif)

## Features

- 🔍 Browse your GitHub repositories with Actions workflows
- 📊 View workflows and their recent runs
- 🔄 Monitor run status in real-time with visual indicators (WIP)
- 👁️ See job status and details for each workflow run (WIP)

## Requirements

### Fonts
gh-ci uses Nerd Font icons for workflow and job status indicators. For the best experience:

1. Install a [Nerd Font](https://www.nerdfonts.com/font-downloads) compatible font of your choice
2. Configure your terminal to use the installed Nerd Font

Without a Nerd Font, the status icons will appear as placeholder characters or missing glyphs.

## Installation

1. Install the `gh` CLI - [instructions](https://github.com/cli/cli?tab=readme-ov-file#installation)
2. Install this extension:

```bash
gh extension install cpaluszek/gh-ci
```

## Configuration
gh-ci automatically creates a default configuration file at first run. You can find or manually edit the config at:

- `$XDG_CONFIG_HOME/gh-ci/config.yaml` (typically `~/.config/gh-ci/config.yaml`)

### Configuration Options
```yaml
github:
  repositories:
    - owner/repo1   # Format: username/repository or organization/repository
    - owner/repo2
```

## Usage

```bash
gh ci
```

Then press `?` for help.

## Development

gh-ci is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions for terminal applications
- [go-gh](https://github.com/cli/go-gh): Go library for the GitHub CLI
- [vhs](https://github.com/charmbracelet/vhs): Gif generation

## Credits

This project was inspired by [gh-dash](https://github.com/dlvhdr/gh-dash) by [dlvhdr](https://github.com/dlvhdr) - a beautiful CLI dashboard for GitHub.
