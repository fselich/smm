package view

import (
	"smm/internal/ui"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/rs/zerolog/log"
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
	log.Info().Msgf("view toast %s", m.text)
	return ui.StyleToast().Width(m.width).Render(m.text)
}
