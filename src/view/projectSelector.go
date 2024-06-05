package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ProjectSelectedMessage struct {
	ProjectId string
}

type ProjectSelector struct {
	teaView textinput.Model
}

func NewProjectSelectorModal() *ProjectSelector {
	projectId := textinput.New()
	projectId.Prompt = "Project ID: "
	projectId.Placeholder = ""
	projectId.Focus()
	projectId.CharLimit = 128
	projectId.Width = 32

	return &ProjectSelector{teaView: projectId}
}

func (p *ProjectSelector) Init() tea.Cmd {
	return nil
}

func (p *ProjectSelector) Value() string {
	return p.teaView.Value()
}

func (p *ProjectSelector) Update(msg tea.Msg) (Modal, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd = func() tea.Msg {
				return ProjectSelectedMessage{p.teaView.Value()}
			}
			cmds = append(cmds, cmd)
		}

	}
	p.teaView, cmd = p.teaView.Update(msg)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *ProjectSelector) View() string {
	return p.teaView.View()
}
