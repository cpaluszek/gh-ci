package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
    Up key.Binding
    Down key.Binding
    Select key.Binding
    Quit key.Binding
    Return key.Binding
    OpenGitHub key.Binding
    Help key.Binding
}

var Keys = &KeyMap{
    Up: key.NewBinding(
        key.WithKeys("up", "k"),
        key.WithHelp("↑/k", "move up"),
        ),
    Down: key.NewBinding(
        key.WithKeys("down", "j"),
        key.WithHelp("↓/j", "move down"),
        ),
    Select: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("enter", "select"),
        ),
    Quit: key.NewBinding(
        key.WithKeys("q", "ctrl+c"),
        key.WithHelp("q", "quit"),
        ),
    Return: key.NewBinding(
        key.WithKeys("esc", "backspace"),
        key.WithHelp("esc", "return"),
        ),
    OpenGitHub: key.NewBinding(
        key.WithKeys("o"),
        key.WithHelp("o", "open in GitHub"),
        ),
    Help: key.NewBinding(
        key.WithKeys("?"),
        key.WithHelp("?", "help"),
        ),
}

func (k KeyMap) ShortHelp() []key.Binding {
    return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding{
        {k.Up, k.Down, k.Select},
        {k.OpenGitHub, k.Return},
        {k.Help, k.Quit},
    }
}
