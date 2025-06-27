package model

import (
	"gcs/gcp"
	"gcs/page"
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
)

type Page interface {
	Init()
	View() string
	Resize(int, int)
	Update(cmd tea.Msg) tea.Cmd
}

type Model struct {
	err       error
	gcp       *gcp.Gcp
	width     int
	height    int
	page      Page
	ProjectId string
}

func New(projectId string) *Model {
	return &Model{ProjectId: projectId}
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
		m.initialize(false)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		return m, nil
	case view.ResizeMessage:
		m.resize()
	}

	if m.page != nil {
		cmd = m.page.Update(msg)
		return m, cmd
	}

	return m, nil
}

type ShowProjectSelectMsg struct {
}

func (m *Model) initialize(showProjectSelect bool) {
	var selected page.CurrentSecret

	m.page = page.NewSecrets(m.gcp, selected.Index())
	m.page.Resize(m.width, m.height)

	if showProjectSelect {
		m.page.Update(tea.KeyMsg{Runes: []rune("p"), Type: tea.KeyRunes})
	}
}

func (m *Model) resize() {
	if m.page != nil && m.width > 0 && m.height > 0 {
		m.page.Resize(m.width, m.height)
	}
}

func (m *Model) View() string {
	if m.page == nil {
		if m.ProjectId != "" {
			m.Update(view.ProjectSelectedMessage{ProjectId: m.ProjectId})
		} else {
			m.initialize(true)
		}

		m.View()
	}

	return m.page.View()
}
