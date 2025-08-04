package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/rs/zerolog/log"
)

type Modal interface {
	View() string
	Init() tea.Cmd
	Update(msg tea.Msg) (Modal, tea.Cmd)
}

type Confirm struct {
	question     string
	message      any
	confirmation *confirmation.Model
}

type ConfirmationResultMessage struct {
	Result bool
	Msg    any
}

func NewConfirm(question string, message any) *Confirm {
	confirm := confirmation.New(question, confirmation.Yes)
	confirm.Template = confirmation.TemplateArrow
	confirm.ResultTemplate = confirmation.ResultTemplateArrow

	confirmModel := confirmation.NewModel(confirm)
	return &Confirm{question: question, message: message, confirmation: confirmModel}
}

func (c *Confirm) View() string {
	return c.confirmation.View()
}

func (c *Confirm) Init() tea.Cmd {
	return c.confirmation.Init()
}

func (c *Confirm) Update(msg tea.Msg) (Modal, tea.Cmd) {
	m, cmd := c.confirmation.Update(msg)

	if cmd != nil {
		message := cmd()
		switch message.(type) {
		case tea.QuitMsg:
			result, _ := m.(*confirmation.Model).Value()
			log.Info().Msgf("Confirmation result: %v", result)
			cmd := func() tea.Msg {
				return ConfirmationResultMessage{result, c.message}
			}
			return c, cmd
		}
	}
	return c, nil
}
