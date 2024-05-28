package view

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Search    key.Binding
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Help      key.Binding
	Edit      key.Binding
	Copy      key.Binding
	Refresh   key.Binding
	ProjectId key.Binding
	Versions  key.Binding
	Quit      key.Binding
}

type Help struct {
	teaView help.Model
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Search, k.Copy, k.Edit, k.Versions, k.Refresh, k.ProjectId, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Edit, k.ProjectId},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),

	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),

	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	ProjectId: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Set ProjectId"),
	),
	Versions: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "View Versions"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
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
	h.teaView.Width = w
}

func (h *Help) SetHeight(w int) {
	h.teaView.Width = w
}

func (h *Help) View() string {
	return h.teaView.View(keys)
}

func (h *Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	var cmd tea.Cmd
	h.teaView, cmd = h.teaView.Update(msg)
	return *h, cmd
}
