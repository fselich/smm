package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/erikgeiser/promptkit/confirmation"
)

type Confirm struct {
	question     string
	message      any
	confirmation *confirmation.Model
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

func (c *Confirm) Update(msg tea.Msg) (*Confirm, tea.Cmd) {
	_, cmd := c.confirmation.Update(msg)
	return c, cmd
}
