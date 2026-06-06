package view

import (
	"smm/internal/ui"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

type DeleteConfirmType int

const (
	DeleteSecret DeleteConfirmType = iota
	DeleteVersion
)

type DeleteConfirmMsg struct {
	FullPath string
	Type     DeleteConfirmType
	Version  int
}

type DeleteConfirm struct {
	title      string
	message    string
	checked    bool
	focusIndex int
	msg        DeleteConfirmMsg
}

const (
	focusCheckbox = iota
	focusCancel
	focusDelete
)

func NewDeleteSecretConfirm(secretName string, fullPath string) *DeleteConfirm {
	return &DeleteConfirm{
		title:   "Delete Secret?",
		message: secretName,
		msg: DeleteConfirmMsg{
			FullPath: fullPath,
			Type:     DeleteSecret,
		},
	}
}

func NewDeleteVersionConfirm(secretName string, fullPath string, version int) *DeleteConfirm {
	return &DeleteConfirm{
		title:   "Delete Version",
		message: secretName,
		msg: DeleteConfirmMsg{
			FullPath: fullPath,
			Type:     DeleteVersion,
			Version:  version,
		},
	}
}

func (c *DeleteConfirm) View() string {
	checkbox := "[ ]"
	if c.checked {
		checkbox = "[x]"
	}

	title := lipgloss.NewStyle().Bold(true).Render(c.title)

	var body strings.Builder
	body.WriteString(title)
	body.WriteString("\n\n  ")
	body.WriteString(c.message)
	body.WriteString("\n\n")

	line := "  " + checkbox + " Yes, delete it"
	if c.focusIndex == focusCheckbox {
		line = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEFA")).Render(line)
	}
	body.WriteString(line)

	body.WriteString("\n\n")

	var cancelBtn, deleteBtn string
	if c.focusIndex == focusCancel {
		cancelBtn = ui.StyleSelectedButton().Render("Cancel")
	} else {
		cancelBtn = ui.StyleUnselectedButton().Render("Cancel")
	}
	if c.checked {
		if c.focusIndex == focusDelete {
			deleteBtn = ui.StyleSelectedButton().Render("Delete")
		} else {
			deleteBtn = ui.StyleUnselectedButton().Render("Delete")
		}
	} else {
		deleteBtn = ui.StyleDisabledButton().Render("Delete")
	}

	body.WriteString("    " + cancelBtn + "  " + deleteBtn)

	return body.String()
}

func (c *DeleteConfirm) Init() tea.Cmd {
	return nil
}

func (c *DeleteConfirm) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			return c, func() tea.Msg {
				return ConfirmationResultMessage{Result: false, Msg: c.msg}
			}
		case "tab", "left", "right":
			if msg.String() == "left" {
				c.focusIndex--
				if c.focusIndex < 0 {
					c.focusIndex = focusDelete
				}
			} else {
				c.focusIndex++
				if c.focusIndex > focusDelete {
					c.focusIndex = focusCheckbox
				}
			}
		case " ", "space":
			if c.focusIndex == focusCheckbox {
				c.checked = !c.checked
			}
		case "enter":
			switch c.focusIndex {
			case focusCheckbox:
				c.checked = !c.checked
			case focusCancel:
				return c, func() tea.Msg {
					return ConfirmationResultMessage{Result: false, Msg: c.msg}
				}
			case focusDelete:
				if c.checked {
					return c, func() tea.Msg {
						return ConfirmationResultMessage{Result: true, Msg: c.msg}
					}
				}
			}
		}
	}
	return c, nil
}
