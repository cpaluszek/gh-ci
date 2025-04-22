# Pipeye

A terminal UI application for monitoring your GitHub Actions workflows and runs.

## Features

- 🔍 Browse your GitHub repositories with Actions workflows
- 📊 View workflows and their recent runs
- 🔄 Monitor run status in real-time with visual indicators
- 👁️ See job status and details for each workflow run

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
  - `q`: Quit

- **Workflow View:**
  - `↑`/`k`, `↓`/`j`: Navigate through workflow runs
  - `Esc`/`Backspace`: Return to repository list

## Development

Pipeye is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions for terminal applications
- [go-github](https://github.com/google/go-github): GitHub API client for Go

### Project Structure

```
pipeye/
├── cmd/pipeye/           # Application entry point
├── internal/
│   ├── app/              # Main application model
│   ├── config/           # Configuration handling
│   ├── github/           # GitHub API interactions
│   ├── models/           # Data models and messages
│   └── ui/               # User interface components
│       ├── render/       # UI rendering functions
│       └── views/        # View components
```

