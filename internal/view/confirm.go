package view

import (
	tea "charm.land/bubbletea/v2"
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
	return c.confirmation.View().Content
}

func (c *Confirm) Init() tea.Cmd {
	return c.confirmation.Init()
}

func (c *Confirm) Update(msg tea.Msg) (Modal, tea.Cmd) {
	m, cmd := c.confirmation.Update(msg)
	c.confirmation = m.(*confirmation.Model)

	if cmd != nil {
		message := cmd()
		if _, ok := message.(tea.QuitMsg); ok {
			result, err := c.confirmation.Value()
			if err != nil {
				log.Error().Err(err).Msg("Confirmation result error")
				return c, nil
			}
			return c, func() tea.Msg {
				return ConfirmationResultMessage{result, c.message}
			}
		}
		return c, cmd
	}
	return c, nil
}
