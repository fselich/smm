package view

import (
	tea "charm.land/bubbletea/v2"
)

type HelpModal struct {
	help   Help
	width  int
	height int
}

func NewHelpModal() *HelpModal {
	help := NewHelp()
	help.teaView.ShowAll = true
	return &HelpModal{
		width:  60,
		height: 20,
		help:   help,
	}
}

func (h *HelpModal) Init() tea.Cmd {
	return nil
}

func (h *HelpModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return h, nil
		}
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
	}
	return h, nil
}

func (h *HelpModal) View() string {
	return h.help.View()
}
