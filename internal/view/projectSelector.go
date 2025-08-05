package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
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
	projectId.SetSuggestions(viper.GetStringSlice("projectids"))
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

			uniqueProjectIds := make(map[string]struct{})
			for _, id := range append(viper.GetStringSlice("projectIds"), p.teaView.Value()) {
				uniqueProjectIds[id] = struct{}{}
			}

			projectIds := make([]string, 0, len(uniqueProjectIds))
			for id := range uniqueProjectIds {
				projectIds = append(projectIds, id)
			}

			viper.Set("selected", p.teaView.Value())
			viper.Set("projectIds", projectIds)
			_ = viper.WriteConfig()

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
