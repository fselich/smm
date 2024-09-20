package view

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type SearchMessage struct {
	Query string
}

type SearchForm struct {
	teaView textinput.Model
}

func NewSearchForm() *SearchForm {
	form := textinput.New()
	form.Prompt = "Query search: "
	form.ShowSuggestions = false
	form.Placeholder = ""
	form.Focus()
	form.CharLimit = 128
	form.Width = 32

	return &SearchForm{teaView: form}
}

func (p *SearchForm) Init() tea.Cmd {
	return nil
}

func (p *SearchForm) Value() string {
	return p.teaView.Value()
}

func (p *SearchForm) Update(msg tea.Msg) (Modal, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			log.Info().Msgf("SearchForm: %v", p.teaView.Value())
			cmd = func() tea.Msg {
				return SearchMessage{p.teaView.Value()}
			}
			cmds = append(cmds, cmd)
		}

	}
	p.teaView, cmd = p.teaView.Update(msg)
	cmds = append(cmds, cmd)
	return p, tea.Batch(cmds...)
}

func (p *SearchForm) View() string {
	return lipgloss.NewStyle().Width(34).Height(1).Render(p.teaView.View())
}
