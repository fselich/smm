package view

import (
	"fmt"
	gcp2 "gcs/gcp"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Version struct {
	title       string
	fullPath    string
	description string
}

func newVersion(title, fullPath, description string) Version {
	return Version{title: title, fullPath: fullPath, description: description}
}

func (v Version) FilterValue() string {
	return v.title
}

func (v Version) Title() string {
	return v.description
}

func (v Version) Description() string {
	return v.description
}

func (v Version) FullPath() string {
	return v.fullPath
}

type SecretVersionsList struct {
	teaView list.Model
}

func NewSecretVersionsList(width, height int, secretName string, title string, gcp *gcp2.Gcp) SecretVersionsList {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)

	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#000000")).
		BorderLeftForeground(lipgloss.Color("#87CEFA")).
		Background(lipgloss.Color("#87CEFA"))
	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Foreground(lipgloss.Color("#87CEFA"))

	myList := list.New([]list.Item{}, d, width, height)
	myList.SetShowHelp(false)
	myList.SetShowTitle(true)
	myList.SetShowStatusBar(false)

	myList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEFA"))

	gcpVersions := gcp.GetSecretVersions(secretName)

	var versionList []list.Item
	for _, version := range gcpVersions {
		versionList = append(versionList, newVersion(version.Name, fmt.Sprintf("%v", version.FullPath), fmt.Sprintf("%v", version.Version)))
	}

	myList.SetItems(versionList)
	myList.Title = title
	return SecretVersionsList{teaView: myList}
}

func (sv *SecretVersionsList) SelectedItem() Version {
	return sv.teaView.SelectedItem().(Version)
}

func (sv *SecretVersionsList) Init() tea.Cmd {
	return nil
}

func (sv *SecretVersionsList) View() string {
	return sv.teaView.View()
}

func (sv *SecretVersionsList) Width() int {
	return sv.teaView.Width()
}

func (sv *SecretVersionsList) Height() int {
	return sv.teaView.Height()
}

func (sv *SecretVersionsList) SetWidth(width int) {
	sv.teaView.SetWidth(width)
}

func (sv *SecretVersionsList) SetHeight(height int) {
	sv.teaView.SetHeight(height)
}

func (sv *SecretVersionsList) Update(msg tea.Msg) (SecretVersionsList, tea.Cmd) {
	cmd := tea.Batch()
	sv.teaView, cmd = sv.teaView.Update(msg)
	return *sv, cmd
}

func (sv *SecretVersionsList) SetTitle(title string) {
	sv.teaView.Title = title
}
