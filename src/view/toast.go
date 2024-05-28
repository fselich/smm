package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

type Toast struct {
	text  string
	width int
	timer time.Time
}

func NewToast() Toast {
	return Toast{text: ""}
}

type ShowToast struct {
	Text string
}

func (m *Toast) SetText(text string) {
	m.text = text
	m.timer = time.Now()
}

func (m *Toast) SetWith(width int) {
	m.width = width
}

func (m *Toast) Update(msg tea.Msg) (Toast, tea.Cmd) {
	switch msg := msg.(type) {
	case ShowToast:
		m.text = msg.Text
		return *m, nil
	}

	if time.Since(m.timer) > 3*time.Second {
		m.text = ""
	}

	return *m, nil
}

func (m *Toast) View() string {
	if time.Since(m.timer) > 3*time.Second {
		m.text = ""
	}
	return lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("#FFFF00")).
		Align(lipgloss.Center).Render(m.text)
}
