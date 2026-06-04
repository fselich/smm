package model

import (
	"smm/internal/client"
	"smm/internal/config"
	"smm/internal/page"
	"smm/internal/view"

	tea "charm.land/bubbletea/v2"
)

type Page interface {
	Init()
	View() string
	Resize(int, int)
	Update(cmd tea.Msg) tea.Cmd
}

type Model struct {
	err       error
	gcp       client.Client
	width     int
	height    int
	page      Page
	ProjectId string
}

func New(projectId string) *Model {
	return &Model{ProjectId: projectId}
}

func (m *Model) Init() tea.Cmd {
	err := m.setProjectId(m.ProjectId)
	m.initialize()

	if m.ProjectId == "" || err != nil {
		m.page.Update(tea.KeyPressMsg{Text: "p"})
	}

	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case view.ProjectSelectedMessage:
		err := m.setProjectId(msg.ProjectId)
		m.initialize()

		if err != nil {
			m.page.Update(view.ShowProjectSelectMsg{TextAlert: "Error setting project ID"})
		}
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

func (m *Model) initialize() {
	var selected page.CurrentSecret

	m.page = page.NewSecrets(m.gcp, selected.Index())
	m.page.Resize(m.width, m.height)
}

func (m *Model) resize() {
	if m.page != nil && m.width > 0 && m.height > 0 {
		m.page.Resize(m.width, m.height)
	}
}

func (m *Model) View() tea.View {
	if m.page == nil {
		return tea.NewView("Error: Application not properly initialized. Please restart.")
	}

	v := tea.NewView(m.page.View())
	v.AltScreen = true
	return v
}

func (m *Model) setProjectId(projectId string) error {
	m.ProjectId = projectId
	if projectId == "" {
		return nil
	}

	projectType := config.GetTypeByProjectId(projectId)
	if projectType == "gcp" {
		gcp, err := client.NewGcp(projectId)
		if err != nil {
			return err
		}
		m.gcp = gcp
	} else {
		fc, err := client.NewFakeClient(projectId)
		if err != nil {
			return err
		}
		m.gcp = fc
	}

	return nil
}
