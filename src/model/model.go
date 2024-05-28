package model

import (
	"gcs/gcp"
	"gcs/page"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
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
	return &Model{status: 0}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case page.ProjectSelectedMsg:
		var err error
		m.gcp, err = gcp.NewGcp(msg.ProjectId)
		if err != nil {
			return m, nil
		}
		m.setStatus(2, msg)
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resize()
		m.page.Resize(msg.Width, msg.Height)
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

func (m *Model) setStatus(status int, msg any) {
	if m.status != 1 && status == 1 {
		m.page = page.NewSelectProject()
		m.page.Resize(m.width, m.height)
		m.status = 1
	}
	if m.status != 2 && status == 2 {
		var selected page.CurrentSecret
		log.Info().Msgf("msg: %v", msg)
		if msg, ok := msg.(page.SetStatusMsg); ok && msg.From == "versions" {
			selected = msg.Data.(page.CurrentSecret)
			log.Info().Msgf("msg is SetStatusMsg and selected is %v", selected)
		}

		m.page = page.NewSecrets(m.gcp, selected.Index())
		m.page.Resize(m.width, m.height)

		m.status = 2
	}

	if m.status != 3 && status == 3 {
		currentSecret := msg.(page.SetStatusMsg).Data.(page.CurrentSecret)
		m.page = page.NewVersions(currentSecret, m.gcp)
		m.status = 3
	}

	m.resize()
}

func (m *Model) resize() {
	if m.width > 0 && m.height > 0 {
		m.page.Resize(m.width, m.height)
	}
}

func (m *Model) View() string {

	switch m.status {
	case 1:
		return m.page.View()
	case 2:
		return m.page.View()
	case 3:
		return m.page.View()
	default:
		m.setStatus(1, nil)
		return "Loading"
	}
}
