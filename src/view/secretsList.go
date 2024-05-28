package view

import (
	gcp2 "gcs/gcp"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"path/filepath"
)

type Secret struct {
	title       string
	fullPath    string
	description string
	index       int
}

func NewSecret(title, fullPath string) Secret {
	return Secret{title: title, fullPath: fullPath, description: "test"}
}

func (t Secret) FilterValue() string {
	return t.title
}

func (t Secret) Title() string {
	return t.title
}

func (t Secret) Description() string {
	return t.description
}

func (t Secret) FullPath() string {
	return t.fullPath
}

func (t Secret) Index() int {
	return t.index
}

func (t Secret) SetIndex(index int) {
	t.index = index
}

type SecretsList struct {
	teaView list.Model
}

func NewSecretsList(width, height int, gcp *gcp2.Gcp) SecretsList {
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
	myList.SetShowTitle(false)
	myList.SetShowStatusBar(false)
	myList.SetShowFilter(false)
	myList.StopSpinner()

	gcpSecrets := gcp.Secrets()

	var secretList []list.Item
	for _, secret := range gcpSecrets {
		secretList = append(secretList, NewSecret(filepath.Base(secret), secret))
	}

	myList.SetItems(secretList)
	return SecretsList{teaView: myList}
}

func (sl *SecretsList) SelectedItem() Secret {
	item := sl.teaView.SelectedItem()

	secret, ok := item.(Secret)
	if !ok {
		return Secret{}
	}

	secret.SetIndex(sl.teaView.Index())
	return secret
}

func (sl *SecretsList) SetItems(myList list.Model, gcpSecrets []string) {
	var secretList []list.Item
	for _, secret := range gcpSecrets {
		secretList = append(secretList, NewSecret(filepath.Base(secret), secret))
	}

	myList.SetItems(secretList)
}

func (sl *SecretsList) Index() int {
	return sl.teaView.Index()
}

func (sl *SecretsList) Init() tea.Cmd {
	return nil
}

func (sl *SecretsList) View() string {
	return sl.teaView.View()
}

func (sl *SecretsList) Width() int {
	return sl.teaView.Width()
}

func (sl *SecretsList) Height() int {
	return sl.teaView.Height()
}

func (sl *SecretsList) SetWidth(width int) {
	sl.teaView.SetWidth(width)
}

func (sl *SecretsList) SetHeight(height int) {
	sl.teaView.SetHeight(height)
}

func (sl *SecretsList) Select(index int) {
	sl.teaView.Select(index)
}

func (sl *SecretsList) Update(msg tea.Msg) (SecretsList, tea.Cmd) {
	var cmd tea.Cmd
	sl.teaView, cmd = sl.teaView.Update(msg)
	return *sl, cmd
}

func (sl *SecretsList) IsFiltering() bool {
	return sl.teaView.SettingFilter()
}

func (sl *SecretsList) IsFiltered() bool {
	return sl.teaView.IsFiltered()
}

func (sl *SecretsList) FilterValue() string {
	return sl.teaView.FilterValue()
}
