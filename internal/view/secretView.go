package view

import (
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type SecretView struct {
	teaView      viewport.Model
	Hidden       bool
	Secret       *Secret
	IsFocused    bool
	IsFiltering  bool
	FilterValue  string
	isFiltered   bool
	content      string
	rawContent   string
	isVisualMode bool
	visualAnchor int
	visualCursor int
	cursorLine   int
}

func NewSecretView(width, height int) SecretView {
	detail := viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
	detail.YPosition = 0

	return SecretView{
		teaView: detail,
		Hidden: false,
		IsFocused: false,
		IsFiltering: false,
		FilterValue: "",
		isFiltered: false,
		cursorLine: -1,
	}
}

func (s *SecretView) SetContent(content string) {
	s.content = content
	s.teaView.SetContent(s.HighlightText())
}

func (s *SecretView) SetRawContent(raw string) {
	s.rawContent = raw
}

func (s *SecretView) View() string {
	s.HighlightText()
	return s.teaView.View()
}

func (s *SecretView) Update(msg tea.Msg) (SecretView, tea.Cmd) {
	cmd := tea.Batch()

	if _, ok := msg.(tea.KeyPressMsg); ok && !s.IsFocused {
		return *s, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && s.isVisualMode {
		switch key.String() {
		case "j", "down":
			lastLine := max(len(strings.Split(strings.TrimRight(s.content, "\n"), "\n"))-1, 0)
			s.visualCursor = min(s.visualCursor+1, lastLine)
			s.updateVisualContent()
			s.scrollIntoView(s.visualCursor)
			return *s, cmd
		case "k", "up":
			s.visualCursor = max(s.visualCursor-1, 0)
			s.updateVisualContent()
			s.scrollIntoView(s.visualCursor)
			return *s, cmd
		case "y":
			selectedText := s.yankSelection()
			s.ExitVisualMode()
			return *s, tea.Batch(
				CopyToClipboard(selectedText),
				func() tea.Msg {
					return ShowToast{Text: "Selection copied to clipboard"}
				},
			)
		case "esc":
			s.ExitVisualMode()
			return *s, cmd
		}

		return *s, cmd
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && s.IsFiltering {
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

	if key, ok := msg.(tea.KeyPressMsg); ok && s.isFiltered {
		switch key.String() {
		case "esc":
			s.IsFiltering = false
			s.isFiltered = false
			s.FilterValue = ""
			return *s, cmd

		case "enter":
			s.teaView.ScrollDown(1)
		}
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			s.IsFiltering = true
		}
	}

	if key, ok := msg.(tea.KeyPressMsg); ok && s.IsFocused && !s.isVisualMode && !s.IsFiltering && !s.isFiltered {
		switch key.String() {
		case "j", "down":
			s.moveCursorDown()
			return *s, cmd
		case "k", "up":
			s.moveCursorUp()
			return *s, cmd
		case "esc":
			if s.cursorLine >= 0 {
				s.cursorLine = -1
				s.teaView.SetContent(s.HighlightText())
				return *s, cmd
			}
		}
	}

	s.teaView, cmd = s.teaView.Update(msg)
	return *s, cmd
}

func (s *SecretView) SetWidth(width int) {
	s.teaView.SetWidth(width)
}

func (s *SecretView) SetHeight(height int) {
	s.teaView.SetHeight(height)
}

func (s *SecretView) Width() int {
	return s.teaView.Width()
}

func (s *SecretView) Height() int {
	return s.teaView.Height()
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
	s.cursorLine = -1
	s.ExitVisualMode()
}

func (s *SecretView) EnterVisualMode() {
	if s.isVisualMode {
		return
	}
	if s.content == "" {
		return
	}
	s.isVisualMode = true
	if s.cursorLine >= 0 {
		s.visualAnchor = s.cursorLine
		s.visualCursor = s.cursorLine
	} else {
		s.visualAnchor = s.teaView.YOffset()
		s.visualCursor = s.teaView.YOffset()
	}
	s.updateVisualContent()
}

func (s *SecretView) ExitVisualMode() {
	if !s.isVisualMode {
		return
	}
	s.cursorLine = s.visualAnchor
	s.isVisualMode = false
	s.visualAnchor = 0
	s.visualCursor = 0
	s.teaView.SetContent(s.HighlightText())
}

func (s *SecretView) IsVisualMode() bool {
	return s.isVisualMode
}

func (s *SecretView) HasCursor() bool {
	return s.cursorLine >= 0
}

func (s *SecretView) ClearCursor() {
	s.cursorLine = -1
	s.teaView.SetContent(s.HighlightText())
}

func (s *SecretView) moveCursorDown() {
	lastLine := max(len(strings.Split(strings.TrimRight(s.content, "\n"), "\n"))-1, 0)
	if s.cursorLine < 0 {
		s.cursorLine = 0
	} else {
		s.cursorLine = min(s.cursorLine+1, lastLine)
	}
	s.updateVisualContent()
	s.scrollIntoView(s.cursorLine)
}

func (s *SecretView) moveCursorUp() {
	if s.cursorLine < 0 {
		s.cursorLine = 0
	} else {
		s.cursorLine = max(s.cursorLine-1, 0)
	}
	s.updateVisualContent()
	s.scrollIntoView(s.cursorLine)
}

func (s *SecretView) renderContent() string {
	if s.isVisualMode {
		lines := strings.Split(s.content, "\n")
		start := min(s.visualAnchor, s.visualCursor)
		end := max(s.visualAnchor, s.visualCursor)
		if start < 0 {
			start = 0
		}
		if end >= len(lines) {
			end = len(lines) - 1
		}
		if start > end {
			return s.content
		}
		for i := start; i <= end; i++ {
			line := strings.ReplaceAll(lines[i], "\033[0m", "\033[0m\033[7m")
			lines[i] = "\033[7m" + line + "\033[27m"
		}
		return strings.Join(lines, "\n")
	}

	if s.cursorLine >= 0 {
		lines := strings.Split(s.content, "\n")
		if s.cursorLine < len(lines) {
			line := strings.ReplaceAll(lines[s.cursorLine], "\033[0m", "\033[0m\033[4m")
			lines[s.cursorLine] = "\033[4m" + line + "\033[24m"
		}
		return strings.Join(lines, "\n")
	}

	return ""
}

func (s *SecretView) updateVisualContent() {
	s.teaView.SetContent(s.renderContent())
}

func (s *SecretView) scrollIntoView(targetLine int) {
	if targetLine < 0 {
		return
	}
	yOffset := s.teaView.YOffset()
	visibleLines := s.teaView.Height()
	if targetLine < yOffset {
		s.teaView.SetYOffset(targetLine)
	} else if targetLine >= yOffset+visibleLines {
		s.teaView.SetYOffset(targetLine - visibleLines + 1)
	}
}

func (s *SecretView) yankSelection() string {
	if s.rawContent == "" {
		return ""
	}
	lines := strings.Split(s.rawContent, "\n")
	start := min(s.visualAnchor, s.visualCursor)
	end := max(s.visualAnchor, s.visualCursor)
	if start < 0 {
		start = 0
	}
	if end >= len(lines) {
		end = len(lines) - 1
	}
	if start > end {
		return ""
	}
	return strings.Join(lines[start:end+1], "\n")
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
