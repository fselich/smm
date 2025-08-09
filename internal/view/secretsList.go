package view

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
	"path/filepath"
	client "smm/internal/client"
	"smm/internal/ui"
	"sort"
	"time"
)

type Secret struct {
	title       string
	fullPath    string
	description string
	index       int
	secretType  string
	version     int
	related     *Secret
	createdAt   time.Time
}

type ResizeMessage struct{}

func NewSecret(title, fullPath, secretType string, version int, createdAt time.Time) Secret {
	return Secret{title: title, fullPath: fullPath, description: "test", secretType: secretType, version: version, createdAt: createdAt}
}

func (t Secret) FilterValue() string {
	if t.secretType == "version" {
		return t.related.title
	}
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

func (t Secret) Type() string {
	return t.secretType
}

func (t Secret) Version() int {
	return t.version
}

func (t Secret) Related() *Secret {
	return t.related
}

func (t Secret) CreatedAt() time.Time {
	return t.createdAt
}

func (t *Secret) SetRelated(secret *Secret) {
	t.related = secret
}

func (t Secret) Hash() string {
	hasher := sha256.New()
	hasher.Write([]byte(t.title))
	return hex.EncodeToString(hasher.Sum(nil)) + ".env"
}

func (t *Secret) SetIndex(index int) {
	t.index = index
}

type SecretsList struct {
	teaView     list.Model
	IsFocused   bool
	SearchQuery string
}

func NewSecretsList(width, height int, gcp client.Client) SecretsList {
	dl := NewListDelegate()
	dl.Styles.SelectedTitle = ui.StyleSelected()
	dl.Styles.NormalTitle = ui.StyleUnselected()

	myList := list.New([]list.Item{}, dl, width, height)
	myList.SetShowHelp(false)
	myList.SetShowTitle(false)
	myList.SetShowStatusBar(false)
	myList.SetShowFilter(false)
	myList.SetShowPagination(true)
	myList.StopSpinner()
	myList.DisableQuitKeybindings()
	myList.Filter = list.UnsortedFilter

	var secretList []list.Item

	if gcp != nil {
		secretInfos, err := gcp.Secrets()
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch secrets")
		} else {
			for _, secretInfo := range secretInfos {
				secretList = append(secretList, NewSecret(secretInfo.Name, secretInfo.FullPath, "current", 0, secretInfo.CreateTime))
			}
		}
	}

	myList.SetItems(secretList)
	return SecretsList{teaView: myList, IsFocused: true}
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
		secretList = append(secretList, NewSecret(filepath.Base(secret), secret, "current", 0, time.Now()))
	}

	myList.SetItems(secretList)
}

func (sl *SecretsList) Index() int {
	return sl.teaView.Index()
}

func (sl *SecretsList) RealIndex() int {
	selected := sl.teaView.SelectedItem()
	for i, item := range sl.teaView.Items() {
		if item == selected {
			return i
		}
	}

	return 0
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

func (sl *SecretsList) SelectByName(name string) {
	sl.teaView.ResetSelected()
	for i, item := range sl.teaView.Items() {
		if item.(Secret).Title() == name {
			sl.teaView.Select(i)
			break
		}
	}
}

func (sl *SecretsList) Update(msg tea.Msg) (SecretsList, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if _, ok := msg.(tea.KeyMsg); ok && !sl.IsFocused {
		return *sl, cmd
	}

	sl.teaView, cmd = sl.teaView.Update(msg)
	cmds = append(cmds, cmd)

	//Fix for resizing after filtering
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "esc":
			resizeCmd := func() tea.Msg {
				return ResizeMessage{}
			}
			cmds = append(cmds, resizeCmd)
		}
	}

	return *sl, tea.Batch(cmds...)
}

func (sl *SecretsList) InsertItem(index int, item Secret) tea.Cmd {
	cmd := sl.teaView.InsertItem(index, item)
	return cmd
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

func (sl *SecretsList) DelVersionItems() bool {
	items := sl.teaView.Items()
	var indexToDelete []int
	deleted := false
	for i, item := range items {
		if item != nil && item.(Secret).Type() == "version" {
			indexToDelete = append(indexToDelete, i)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(indexToDelete)))
	for _, index := range indexToDelete {
		sl.teaView.RemoveItem(index)
		deleted = true
	}

	return deleted
}

func (sl *SecretsList) ResetFilter() {
	sl.teaView.ResetFilter()
}

func (sl *SecretsList) DeepSearch(query string, gcp client.Client) {

	var secretList []list.Item

	if query == "" {
		secretInfos, err := gcp.Secrets()
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch secrets")
			return
		}
		for _, secretInfo := range secretInfos {
			secretList = append(secretList, NewSecret(secretInfo.Name, secretInfo.FullPath, "current", 0, secretInfo.CreateTime))
		}
	} else {
		searchResults, err := gcp.SearchInSecrets(query)
		if err != nil {
			log.Error().Err(err).Msg("failed to search secrets")
			return
		}
		for _, secretInfo := range searchResults {
			secretList = append(secretList, NewSecret(secretInfo.Name, secretInfo.FullPath, "current", 0, secretInfo.CreateTime))
		}
	}

	sl.SearchQuery = query

	sl.teaView.SetItems(secretList)
}

func (sl *SecretsList) ToggleFocus() {
	if sl.IsFocused {
		sl.IsFocused = false
	} else {
		sl.IsFocused = true
	}
}
