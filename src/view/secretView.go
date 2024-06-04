package view

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type SecretView struct {
	teaView viewport.Model
	Hidden  bool
	Secret  *Secret
}

func NewSecretView(width, height int) SecretView {
	detail := viewport.New(width, height)
	detail.HighPerformanceRendering = false
	detail.YPosition = 0

	return SecretView{teaView: detail, Hidden: false}
}

func (s *SecretView) SetContent(content string) {
	s.teaView.SetContent(content)
}

func (s *SecretView) View() string {
	return s.teaView.View()
}

func (s *SecretView) Update(msg tea.Msg) (SecretView, tea.Cmd) {
	cmd := tea.Batch()
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
