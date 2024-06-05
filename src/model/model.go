package model

import (
	"gcs/gcp"
	"gcs/page"
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	err    error
	gcp    *gcp.Gcp
	status int
	width  int
	height int
	page   page.Page
}

func New() *Model {
	return &Model{status: 1}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case view.ProjectSelectedMessage:
		var err error
		m.gcp, err = gcp.NewGcp(msg.ProjectId)
		if err != nil {
			return m, nil
		}
		m.setStatus(2, nil)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil
	case page.SetStatusMsg:
		m.setStatus(msg.Status, msg)
		return m, nil
	}

	if m.page != nil {
		cmd = m.page.Update(msg)
		return m, cmd
	}

	return m, nil
}

type ShowProjectSelectMsg struct {
}

func (m *Model) setStatus(status int, msg any) {
	if status == 2 {
		var selected page.CurrentSecret

		m.page = page.NewSecrets(m.gcp, selected.Index())
		m.page.Resize(m.width, m.height)
		if _, ok := msg.(ShowProjectSelectMsg); ok {
			m.page.Update(tea.KeyMsg{Runes: []rune("p"), Type: tea.KeyRunes})
		}
		m.status = 2
	}

	m.resize()
}

func (m *Model) resize() {
	if m.page != nil && m.width > 0 && m.height > 0 {
		m.page.Resize(m.width, m.height)
	}
}

func (m *Model) View() string {
	if m.page == nil {
		m.setStatus(2, ShowProjectSelectMsg{})
		m.View()
	}

	return m.page.View()
}
