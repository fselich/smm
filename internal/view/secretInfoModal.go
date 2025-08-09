package view

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"smm/internal/client"
)

type SecretInfoModal struct {
	secretInfo client.SecretInfo
	width      int
	height     int
}

type SecretInfoDisplayMsg struct {
	SecretInfo client.SecretInfo
}

func NewSecretInfoModal(secretInfo client.SecretInfo) *SecretInfoModal {
	return &SecretInfoModal{
		secretInfo: secretInfo,
		width:      60,
		height:     20,
	}
}

func (s *SecretInfoModal) Init() tea.Cmd {
	return nil
}

func (s *SecretInfoModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, nil
		}
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}
	return s, nil
}

func (s *SecretInfoModal) View() string {
	var sections []string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#87CEFA")).
		MarginBottom(1)

	sections = append(sections, titleStyle.Render("Secret Information"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0"))

	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Name: "),
			valueStyle.Render(s.secretInfo.Name),
		),
	)

	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Full Path: "),
			valueStyle.Render(s.secretInfo.FullPath),
		),
	)

	createdTime := s.secretInfo.CreateTime.Format("2006-01-02 15:04:05 UTC")
	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Created: "),
			valueStyle.Render(createdTime),
		),
	)

	age := time.Since(s.secretInfo.CreateTime)
	var ageStr string
	if age.Hours() > 24 {
		days := int(age.Hours() / 24)
		if days == 1 {
			ageStr = "1 day ago"
		} else {
			ageStr = fmt.Sprintf("%d days ago", days)
		}
	} else if age.Hours() > 1 {
		hours := int(age.Hours())
		if hours == 1 {
			ageStr = "1 hour ago"
		} else {
			ageStr = fmt.Sprintf("%d hours ago", hours)
		}
	} else {
		minutes := int(age.Minutes())
		if minutes <= 1 {
			ageStr = "just now"
		} else {
			ageStr = fmt.Sprintf("%d minutes ago", minutes)
		}
	}

	sections = append(sections,
		lipgloss.JoinHorizontal(lipgloss.Left,
			labelStyle.Render("Age: "),
			valueStyle.Render(ageStr),
		),
	)

	if len(s.secretInfo.Labels) > 0 {
		sections = append(sections, "") // Empty line
		sections = append(sections, labelStyle.Render("Labels:"))

		for key, value := range s.secretInfo.Labels {
			labelEntry := fmt.Sprintf("  %s: %s", key, value)
			sections = append(sections, valueStyle.Render(labelEntry))
		}
	}

	if len(s.secretInfo.Annotations) > 0 {
		sections = append(sections, "")
		sections = append(sections, labelStyle.Render("Annotations:"))

		for key, value := range s.secretInfo.Annotations {
			displayValue := value
			if len(displayValue) > 50 {
				displayValue = displayValue[:47] + "..."
			}
			annotationEntry := fmt.Sprintf("  %s: %s", key, displayValue)
			sections = append(sections, valueStyle.Render(annotationEntry))
		}
	}

	sections = append(sections, "")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5a5a5a")).
		Italic(true)
	sections = append(sections, footerStyle.Render("Press ESC to close"))

	content := strings.Join(sections, "\n")

	return content
}
