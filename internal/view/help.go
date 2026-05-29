package view

import (
	"smm/internal/config"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type keyMap struct {
	Filter     key.Binding
	Search     key.Binding
	NewSecret  key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Help       key.Binding
	NewVersion key.Binding
	Copy       key.Binding
	Refresh    key.Binding
	Restore    key.Binding
	ProjectId  key.Binding
	Versions   key.Binding
	Visual     key.Binding
	Yank       key.Binding
	More       key.Binding
	Info       key.Binding
	Quit       key.Binding
}

type Help struct {
	teaView help.Model
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	if config.ExpermientalEnabled() {
		return []key.Binding{k.Filter, k.Search, k.NewSecret, k.NewVersion, k.Versions, k.Restore, k.Copy, k.Visual, k.Yank, k.More, k.Quit}
	} else {
		return []key.Binding{k.Filter, k.Search, k.NewSecret, k.NewVersion, k.Versions, k.Restore, k.Copy, k.More, k.Quit}
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	if config.ExpermientalEnabled() {
		return [][]key.Binding{
			{k.Up, k.Down, k.Left, k.Right, k.Filter, k.Search, k.Refresh},
			{k.NewSecret, k.NewVersion, k.Versions, k.Restore, k.Info, k.Copy},
			{k.Visual, k.Yank, k.ProjectId, k.Help, k.More, k.Quit},
		}
	} else {
		return [][]key.Binding{
			{k.Up, k.Down, k.Left, k.Right, k.Filter, k.Search, k.Refresh},
			{k.NewSecret, k.NewVersion, k.Versions, k.Restore, k.Info, k.Copy},
			{k.ProjectId, k.Help, k.More, k.Quit},
		}
	}
}

var keys = keyMap{
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter"),
	),
	Search: key.NewBinding(
		key.WithKeys("ctr+f"),
		key.WithHelp("ctr+f", "Search"),
	),
	NewSecret: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "New secret"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "Move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "Move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "Move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "Move right"),
	),

	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Help"),
	),

	NewVersion: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New version"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Copy"),
	),
	Restore: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Restore"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Refresh"),
	),
	ProjectId: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Set ProjectId"),
	),
	Versions: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "View versions"),
	),
	Visual: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "Visual"),
	),
	Yank: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "Yank"),
	),
	More: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "More"),
	),
	Info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "Secret info"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "Quit"),
	),
}

func NewHelp() Help {
	h := help.New()
	h.ShowAll = false
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#C0C0C0"))
	h.Styles.FullKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#C0C0C0"))
	h.Styles.FullDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#C0C0C0"))

	return Help{teaView: h}
}

func (h *Help) SetWidth(w int) {
	h.teaView.SetWidth(w)
}

func (h *Help) SetHeight(w int) {
	h.teaView.SetWidth(w)
}

func (h *Help) View() string {
	return h.teaView.View(keys)
}

func (h *Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	var cmd tea.Cmd
	h.teaView, cmd = h.teaView.Update(msg)
	return *h, cmd
}
