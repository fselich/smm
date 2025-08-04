package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"smm/internal/gcp"
	"smm/internal/page"
	"smm/internal/view"
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
	m.setProjectId(m.ProjectId)
	m.initialize()

	if m.ProjectId == "" {
		m.page.Update(tea.KeyMsg{Runes: []rune("p"), Type: tea.KeyRunes})
	}

	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case view.ProjectSelectedMessage:
		m.setProjectId(msg.ProjectId)
		m.initialize()
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

func (m *Model) View() string {
	if m.page == nil {
		m.View()
	}

	return m.page.View()
}

func (m *Model) setProjectId(projectId string) {
	m.ProjectId = projectId
	if projectId == "" {
		return
	}

	var err error
	m.gcp, err = gcp.NewGcp(projectId)
	if err != nil {
		log.Fatal().Msgf("Error initializing Gcp: %v", err)
	}
}
