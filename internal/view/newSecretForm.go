package view

import (
	"fmt"
	"smm/internal/braceexpand"
	"smm/internal/ui"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/rs/zerolog/log"
)

type CreateSecretsMessage struct {
	Names []string
}

type NewSecretForm struct {
	teaView     textinput.Model
	parsedNames []string
	parseErr    error
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

	p.teaView, cmd = p.teaView.Update(msg)
	cmds = append(cmds, cmd)

	value := p.teaView.Value()
	if strings.Contains(value, "{") {
		p.parsedNames, p.parseErr = braceexpand.ExpandBraces(value)
	} else if value != "" {
		p.parsedNames = []string{value}
		p.parseErr = nil
	} else {
		p.parsedNames = nil
		p.parseErr = nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			if value == "" {
				return p, nil
			}
			if p.parseErr != nil {
				log.Info().Err(p.parseErr).Msg("NewSecretForm: invalid pattern")
				return p, nil
			}
			names := p.parsedNames
			if len(names) == 0 {
				return p, nil
			}
			log.Info().Msgf("NewSecretForm: %v", names)
			cmd = func() tea.Msg {
				return CreateSecretsMessage{Names: names}
			}
			cmds = append(cmds, cmd)
		}
	}

	return p, tea.Batch(cmds...)
}

func (p *NewSecretForm) View() string {
	view := p.teaView.View()
	preview := p.previewText()
	if preview != "" {
		view += "\n\n" + ui.StyleLow().Render(preview)
	}
	return lipgloss.NewStyle().Width(48).Height(3).Render(view)
}

func (p *NewSecretForm) previewText() string {
	if p.parseErr != nil {
		if strings.Contains(p.parseErr.Error(), "exceeds limit") {
			return "Exceeds limit"
		}
		return "Invalid pattern"
	}
	if len(p.parsedNames) <= 1 {
		return ""
	}
	return fmt.Sprintf("%d secrets will be created", len(p.parsedNames))
}
