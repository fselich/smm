package view

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type SecretView struct {
	teaView     viewport.Model
	Hidden      bool
	Secret      *Secret
	IsFocused   bool
	IsFiltering bool
	FilterValue string
	isFiltered  bool
	content     string
}

func NewSecretView(width, height int) SecretView {
	detail := viewport.New(width, height)
	detail.HighPerformanceRendering = false
	detail.YPosition = 0

	return SecretView{teaView: detail, Hidden: false, IsFocused: false, IsFiltering: false, FilterValue: "", isFiltered: false}
}

func (s *SecretView) SetContent(content string) {
	s.content = content
	s.teaView.SetContent(s.HighlightText())
}

func (s *SecretView) View() string {
	s.HighlightText()
	return s.teaView.View()
}

func (s *SecretView) Update(msg tea.Msg) (SecretView, tea.Cmd) {
	cmd := tea.Batch()

	if _, ok := msg.(tea.KeyMsg); ok && !s.IsFocused {
		return *s, cmd
	}

	if key, ok := msg.(tea.KeyMsg); ok && s.IsFiltering {
		switch key.String() {
		case "esc":
			s.IsFiltering = false
			s.isFiltered = false
			s.FilterValue = ""
		case "backspace":
			if len(s.FilterValue) > 0 {
				s.FilterValue = s.FilterValue[:len(s.FilterValue)-1]
			}
		case "enter":
			s.IsFiltering = false
			s.isFiltered = true
		default:
			s.FilterValue = s.FilterValue + key.String()
		}

		return *s, cmd
	}

	if key, ok := msg.(tea.KeyMsg); ok && s.isFiltered {
		switch key.String() {
		case "esc":
			s.IsFiltering = false
			s.isFiltered = false
			s.FilterValue = ""
			return *s, cmd

		case "enter":
			s.teaView.LineDown(1)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			s.IsFiltering = true
		}
	}

	s.teaView, cmd = s.teaView.Update(msg)
	return *s, cmd
}

func (s *SecretView) SetWidth(width int) {
	s.teaView.Width = width
}

func (s *SecretView) SetHeight(height int) {
	s.teaView.Height = height
}

func (s *SecretView) Width() int {
	return s.teaView.Width
}

func (s *SecretView) Height() int {
	return s.teaView.Height
}

func (s *SecretView) ToggleFocus() {
	if s.IsFocused {
		s.IsFocused = false
	} else {
		s.IsFocused = true
	}
	s.IsFiltering = false
	s.isFiltered = false
	s.FilterValue = ""
}

func (s *SecretView) HighlightText() string {
	if s.FilterValue == "" || !s.isFiltered {
		return s.content
	}

	return strings.ReplaceAll(s.content, s.FilterValue, "\033[1;37;41m"+s.FilterValue+"\033[0m")
}

func (s *SecretView) SetFilteredValue(value string) {
	if value == "" {
		s.isFiltered = false
		s.FilterValue = ""
	} else {
		s.isFiltered = true
		s.FilterValue = value
	}
}
