package keys

import "github.com/charmbracelet/bubbles/key"

// TODO: implement refresh on all sections
type KeyMap struct {
	Up         key.Binding
	Down       key.Binding
	Select     key.Binding
	Quit       key.Binding
	Return     key.Binding
	OpenGitHub key.Binding
	Tab        key.Binding
	ShiftTab   key.Binding
	Help       key.Binding
	Refresh    key.Binding
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
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch to table"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "switch back table"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "return"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
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
		{k.Tab, k.ShiftTab},
		{k.Refresh},
		{k.Help, k.Quit},
	}
}
