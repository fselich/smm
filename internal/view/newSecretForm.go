package view

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"
)

type CreateSecretMessage struct {
	Name string
}

type NewSecretForm struct {
	teaView textinput.Model
}

func NewSecretFormModal() *NewSecretForm {
	form := textinput.New()
	form.Prompt = "Secret name: "
	form.ShowSuggestions = false
	form.Placeholder = ""
	form.Focus()
	form.CharLimit = 128
	form.SetWidth(32)

	return &NewSecretForm{teaView: form}
}

func (p *NewSecretForm) Init() tea.Cmd {
	return nil
}

func (p *NewSecretForm) Update(msg tea.Msg) (Modal, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if p.teaView.Value() == "" {
				return p, nil
			}
			log.Info().Msgf("NewSecretForm: %v", p.teaView.Value())
			cmd = func() tea.Msg {
				return CreateSecretMessage{p.teaView.Value()}
			}
			cmds = append(cmds, cmd)
		}

	}
	p.teaView, cmd = p.teaView.Update(msg)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *NewSecretForm) View() string {
	return lipgloss.NewStyle().Width(34).Height(1).Render(p.teaView.View())
}
