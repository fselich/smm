package view

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"smm/internal/client"
)

// Styles for the SecretInfoModal
type secretInfoStyles struct {
	title  lipgloss.Style
	label  lipgloss.Style
	value  lipgloss.Style
	footer lipgloss.Style
}

// newSecretInfoStyles creates and returns the styling configuration
func newSecretInfoStyles() secretInfoStyles {
	return secretInfoStyles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#87CEFA")).
			MarginBottom(1),
		label: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")),
		value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C0C0C0")),
		footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5a5a5a")).
			Italic(true),
	}
}

type SecretInfoModal struct {
	secretInfo   client.SecretInfo
	selectedItem Secret
	width        int
	height       int
}

type SecretInfoDisplayMsg struct {
	SecretInfo client.SecretInfo
}

func NewSecretInfoModal(secretInfo client.SecretInfo, selectedItem Secret) *SecretInfoModal {
	return &SecretInfoModal{
		secretInfo:   secretInfo,
		selectedItem: selectedItem,
		width:        60,
		height:       20,
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
	styles := newSecretInfoStyles()
	var sections []string

	// Build title section
	sections = append(sections, s.buildTitle(styles))

	// Build secret info section
	sections = append(sections, "")
	sections = append(sections, s.buildSecretInfoSection(styles)...)

	// Build version info section
	sections = append(sections, "")
	sections = append(sections, s.buildVersionInfoSection(styles)...)

	// Build footer
	sections = append(sections, "")
	sections = append(sections, styles.footer.Render("Press ESC to close"))

	return strings.Join(sections, "\n")
}

// buildTitle creates the modal title based on the selected item type
func (s *SecretInfoModal) buildTitle(styles secretInfoStyles) string {
	var title string
	if s.selectedItem.Type() == "version" {
		title = fmt.Sprintf("Secret Information - Version %d", s.selectedItem.Version())
	} else {
		title = "Secret Information - Current Version"
	}
	return styles.title.Render(title)
}

// buildSecretInfoSection creates the general secret information section
func (s *SecretInfoModal) buildSecretInfoSection(styles secretInfoStyles) []string {
	var sections []string

	sections = append(sections, styles.label.Render("Secret Info:"))
	sections = append(sections, s.buildFieldRow(styles, "  Name: ", s.secretInfo.Name))
	sections = append(sections, s.buildFieldRow(styles, "  Full Path: ", s.secretInfo.FullPath))

	createdTime := s.secretInfo.CreateTime.Format("2006-01-02 15:04:05 UTC")
	sections = append(sections, s.buildFieldRow(styles, "  Created: ", createdTime))

	ageStr := formatTimeAge(time.Since(s.secretInfo.CreateTime))
	sections = append(sections, s.buildFieldRow(styles, "  Age: ", ageStr))

	// Add labels if present
	if len(s.secretInfo.Labels) > 0 {
		sections = append(sections, s.buildMapSection(styles, "  Labels:", s.secretInfo.Labels, false)...)
	}

	// Add annotations if present
	if len(s.secretInfo.Annotations) > 0 {
		sections = append(sections, s.buildMapSection(styles, "  Annotations:", s.secretInfo.Annotations, true)...)
	}

	return sections
}

// buildVersionInfoSection creates the version-specific information section
func (s *SecretInfoModal) buildVersionInfoSection(styles secretInfoStyles) []string {
	var sections []string

	sections = append(sections, styles.label.Render("Selected Version Info:"))

	// Version number or "current"
	versionStr := "current"
	if s.selectedItem.Type() != "current" {
		versionStr = fmt.Sprintf("%d", s.selectedItem.Version())
	}
	sections = append(sections, s.buildFieldRow(styles, "  Version: ", versionStr))

	// Version creation time and age
	createdTime := s.selectedItem.CreatedAt().Format("2006-01-02 15:04:05 UTC")
	sections = append(sections, s.buildFieldRow(styles, "  Created: ", createdTime))

	versionAgeStr := formatTimeAge(time.Since(s.selectedItem.CreatedAt()))
	sections = append(sections, s.buildFieldRow(styles, "  Age: ", versionAgeStr))

	return sections
}

// buildFieldRow creates a label-value row
func (s *SecretInfoModal) buildFieldRow(styles secretInfoStyles, label, value string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left,
		styles.label.Render(label),
		styles.value.Render(value),
	)
}

// buildMapSection creates a section for displaying key-value maps (labels/annotations)
func (s *SecretInfoModal) buildMapSection(styles secretInfoStyles, header string, data map[string]string, truncateValues bool) []string {
	var sections []string
	sections = append(sections, styles.label.Render(header))

	for key, value := range data {
		displayValue := value
		if truncateValues && len(displayValue) > 50 {
			displayValue = displayValue[:47] + "..."
		}
		entry := fmt.Sprintf("    %s: %s", key, displayValue)
		sections = append(sections, styles.value.Render(entry))
	}

	return sections
}

// formatTimeAge converts a time.Duration to a human-readable age string
func formatTimeAge(age time.Duration) string {
	if age.Hours() > 24 {
		days := int(age.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	if age.Hours() > 1 {
		hours := int(age.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}

	minutes := int(age.Minutes())
	if minutes <= 1 {
		return "just now"
	}
	return fmt.Sprintf("%d minutes ago", minutes)
}
