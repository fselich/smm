package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"smm/internal/config"
)

type ProjectSelectedMessage struct {
	ProjectId string
}

type ShowProjectSelectMsg struct {
	TextAlert string
}

type ProjectSelector struct {
	teaView    textinput.Model
	alertText  string
	alertStyle lipgloss.Style
}

func NewProjectSelectorModal() *ProjectSelector {
	projectId := textinput.New()
	projectId.Prompt = "Project ID: "
	projectId.ShowSuggestions = true
	projectId.SetSuggestions(config.GetProjectIDs())
	projectId.Placeholder = ""
	projectId.Focus()
	projectId.CharLimit = 128
	projectId.Width = 32

	alertStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	return &ProjectSelector{teaView: projectId, alertStyle: alertStyle}
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

			config.SetSelectedProject(p.teaView.Value())
			config.AddProjectID(p.teaView.Value())
			_ = config.Save()

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

func (p *ProjectSelector) SetAlert(text string) {
	p.alertText = text
}

func (p *ProjectSelector) ClearAlert() {
	p.alertText = ""
}

func (p *ProjectSelector) View() string {
	view := p.teaView.View()

	if p.alertText != "" {
		alertView := p.alertStyle.
			Width(32).
			Align(lipgloss.Center).
			Render(p.alertText)
		view = lipgloss.JoinVertical(lipgloss.Left, view, alertView)
	}

	return lipgloss.NewStyle().Width(34).Render(view)
}
