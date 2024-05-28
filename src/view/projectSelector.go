package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectSelectedMsg struct {
	projectId string
}

type ProjectSelector struct {
	teaView textinput.Model
}

func NewProjectSelector() ProjectSelector {
	projectId := textinput.New()
	projectId.Placeholder = "Project ID"
	projectId.Focus()
	projectId.CharLimit = 128
	projectId.Width = 32

	return ProjectSelector{teaView: projectId}
}

func (p *ProjectSelector) Value() string {
	return p.teaView.Value()
}

func (p *ProjectSelector) Update(msg tea.Msg) (ProjectSelector, tea.Cmd) {
	var cmd tea.Cmd
	p.teaView, cmd = p.teaView.Update(msg)
	return *p, cmd
}

func (p *ProjectSelector) View() string {
	return p.teaView.View()
}
