package page

import (
	"fmt"
	"os"
	"path/filepath"
	"smm/internal/client"
	"smm/internal/editor"
	"smm/internal/ui"
	"smm/internal/view"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"github.com/tiagomelo/go-clipboard/clipboard"
)

type Secrets struct {
	gcp        client.Client
	components secretsComponents
	Modal      view.Modal
	ListWidth  int
}

type CurrentSecret struct {
	name  string
	title string
	index int
}

type secretsComponents struct {
	list   *view.SecretsList
	detail *view.SecretView
	help   *view.Help
	toast  *view.Toast
}

type RestoreSecretMsg struct {
	FullPath string
	Title    string
	Version  int
}

func (S CurrentSecret) Name() string {
	return S.name
}

func (S CurrentSecret) Title() string {
	return S.title
}

func (S CurrentSecret) Index() int {
	return S.index
}

func (s *Secrets) View() string {
	borderedList := ui.StyleBorder(s.components.list.IsFocused).
		Width(s.components.list.Width()).
		Render(s.components.list.View())
	borderedDetail := ui.StyleBorder(s.components.detail.IsFocused).
		Render(s.components.detail.View())

	borderedHelp := ui.StyleLowBorder().
		Width(s.components.list.Width() + s.components.detail.Width() + 2).
		Render(s.components.help.View())

	var x int
	if s.components.list.IsFiltering() {
		x = ((s.components.list.Width() - len(s.components.list.FilterValue())) / 2) - 1
		listTitle := ui.StyleBorderTitle().Render(s.components.list.FilterValue())
		borderedList = ui.PlaceOverlay(x, 0, "\""+listTitle+"█\"", borderedList, false)
	} else if s.components.list.IsFiltered() {
		x = (s.components.list.Width() - len(s.components.list.FilterValue())) / 2
		listTitle := ui.StyleBorderTitle().Render(s.components.list.FilterValue())
		borderedList = ui.PlaceOverlay(x, 0, listTitle, borderedList, false)
	}

	if s.components.detail.IsFiltering {
		x = (s.components.detail.Width() - len(s.components.detail.FilterValue)) / 2
		detailTitle := ui.StyleBorderTitle().Render(s.components.detail.FilterValue)
		borderedDetail = ui.PlaceOverlay(x, 0, detailTitle, borderedDetail, false)
	} else {
		x = (s.components.detail.Width() - len(s.components.list.SelectedItem().Title())) / 2
		detailTitle := ui.StyleBorderTitle().Render(s.components.list.SelectedItem().Title())
		borderedDetail = ui.PlaceOverlay(x, 0, detailTitle, borderedDetail, false)
	}

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, borderedList, borderedDetail),
		lipgloss.JoinHorizontal(lipgloss.Bottom, borderedHelp),
		lipgloss.JoinHorizontal(lipgloss.Bottom, s.components.toast.View()),
	)

	if s.Modal != nil {
		render = ui.ModalOverlay(s.Modal.View(), render)
	}

	return render
}

func (s *Secrets) Resize(width int, height int) {

	if s.components.list.IsFiltering() {
		s.components.list.SetHeight(height - 12)
	} else {
		s.components.list.SetHeight(height - 6)
	}

	s.components.list.SetWidth(s.ListWidth)
	s.components.detail.SetWidth(width - 5 - s.components.list.Width())
	s.components.detail.SetHeight(height - 6)
	s.components.help.SetWidth(s.components.list.Width() + s.components.detail.Width() + 2)
	s.components.toast.SetWith(width)
}

func (s *Secrets) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case view.ShowToast:
		s.components.toast.SetText(msg.Text)
		return nil
	case view.ShowProjectSelectMsg:
		modal := view.NewProjectSelectorModal()
		modal.SetAlert(msg.TextAlert)
		s.Modal = modal
		s.Modal.Init()
	case view.SearchMessage:
		s.components.list.DeepSearch(msg.Query, s.gcp)
		s.Modal = nil
		s.components.detail.SetFilteredValue(msg.Query)
	case view.ConfirmationResultMessage:
		switch msg.Msg.(type) {
		case editor.EditFinishedMsg:
			s.Modal = nil
			newVersionMessage := msg.Msg.(editor.EditFinishedMsg)
			log.Info().Msgf("Confirmation result in secrets: %v", msg.Result)
			if msg.Result {
				var title string
				if newVersionMessage.CurrentSecret.Type() == "version" {
					title = newVersionMessage.CurrentSecret.Related().Title()
				} else {
					title = newVersionMessage.CurrentSecret.Title()
				}
				log.Info().Msg("Creating new secret")
				log.Info().Msgf("Creating new secret based on %v", title)
				err := s.gcp.AddSecretVersion(title, newVersionMessage.SecretData)
				if err != nil {
					log.Error().Msgf("Error creating new secret: %v", err)
				}
			}
		case RestoreSecretMsg:
			s.components.toast.SetText("Restoring version")
			s.Modal = nil
			restoreMessage := msg.Msg.(RestoreSecretMsg)
			if msg.Result {
				log.Info().Msgf("Restoring secret %v version %v", restoreMessage.Title, restoreMessage.Version)
				secretData, err := s.gcp.GetSecretVersion(restoreMessage.FullPath, strconv.Itoa(restoreMessage.Version))
				if err != nil {
					log.Error().Err(err).Msg("Error getting secret version for restore")
					s.components.toast.SetText("Error getting secret version")
					return nil
				}
				log.Info().Msg("Restoring secret")
				err = s.gcp.AddSecretVersion(restoreMessage.Title, secretData)
				if err != nil {
					log.Error().Msgf("Error creating new secret: %v", err)
					s.components.toast.SetText("Error restoring secret")
					return nil
				}
				s.components.toast.SetText("Secret restored")
			} else {
				s.components.toast.SetText("Restore canceled")
			}
		}
	}

	if s.Modal == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if s.components.list.IsFiltering() == false && s.components.detail.IsFiltering == false {
				switch msg.String() {
				case "n":
					tempDir := os.TempDir()
					hash := s.components.list.SelectedItem().Hash()
					filename := filepath.Join(tempDir, hash)
					f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
					if err != nil {
						s.components.toast.SetText("Failed to create temporary file")
						return nil
					}
					defer f.Close()

					secretName := s.components.list.SelectedItem().FullPath()

					var secretData string
					if s.components.list.SelectedItem().Type() == "version" {
						version := s.components.list.SelectedItem().Version()
						data, err := s.gcp.GetSecretVersion(secretName, strconv.Itoa(version))
						if err != nil {
							log.Error().Err(err).Msg("Error getting secret version")
							return nil
						}
						secretData = string(data)
					} else {
						data, err := s.gcp.GetSecret(secretName)
						if err != nil {
							log.Error().Err(err).Msg("Error getting secret")
							return nil
						}
						secretData = string(data)
					}

					_, err = f.WriteString(secretData)
					if err != nil {
						s.components.toast.SetText("Failed to write secret to temporary file")
						return nil
					}

					return editor.OpenEditor(secretData, s.components.list.SelectedItem())
				case "r":
					if s.components.list.SelectedItem().Type() == "current" {
						s.components.toast.SetText("Cannot restore current version")
						return nil
					}
					s.components.toast.SetText(fmt.Sprintf("Restoring version"))
					msg := RestoreSecretMsg{
						FullPath: s.components.list.SelectedItem().Related().FullPath(),
						Title:    s.components.list.SelectedItem().Related().Title(),
						Version:  s.components.list.SelectedItem().Version(),
					}
					s.Modal = view.NewConfirm("Do you want to restore this secret version?", msg)
					s.Modal.Init()
				case "p":
					s.Modal = view.NewProjectSelectorModal()
					s.Modal.Init()
				case "ctrl+c":
					return tea.Quit
				case "v":
					selected := s.components.list.SelectedItem()
					if selected.Type() == "current" {
						deleted := s.components.list.DelVersionItems()
						if !deleted {
							versions, err := s.gcp.GetSecretVersions(selected.FullPath())
							if err != nil {
								log.Error().Err(err).Msg("Error getting secret versions")
								s.components.toast.SetText("Error getting secret versions")
								return nil
							}
							versions = versions[1:]
							s.components.toast.SetText(fmt.Sprintf("Secret has %v versions", len(versions)))
							for i, version := range versions {
								secret := view.NewSecret(strconv.Itoa(version.Version), version.FullPath, "version", version.Version, version.CreatedAt)
								secret.SetRelated(&selected)
								cmd = s.components.list.InsertItem(s.components.list.RealIndex()+1+i, secret)
							}
						}
					}
					s.components.list.Select(selected.Index())
					return cmd
				case "c":
					secretName := s.components.list.SelectedItem().FullPath()
					data, err := s.gcp.GetSecret(secretName)
					if err != nil {
						log.Error().Err(err).Msg("Error getting secret for clipboard")
						s.components.toast.SetText("Error getting secret")
					} else {
						secretData := string(data)
						c := clipboard.New()
						err = c.CopyText(secretData)
						if err != nil {
							log.Error().Err(err).Msg("Error copying to clipboard")
							s.components.toast.SetText("Failed to copy to clipboard")
						} else {
							s.components.toast.SetText("Secret copied to clipboard")
						}
					}
				case "?":
					s.Modal = view.NewProjectSelectorModal()
					s.Modal.Init()
				case "esc":
					s.Init()
					resizeCmd := func() tea.Msg {
						return view.ResizeMessage{}
					}
					cmds = append(cmds, resizeCmd)
					s.components.toast.SetText("Secrets refreshed")
					s.components.list.SearchQuery = ""
					s.components.detail.SetFilteredValue("")
				case "ctrl+f":
					s.Modal = view.NewSearchForm()
					s.Modal.Init()
				case "i":
					selected := s.components.list.SelectedItem()
					if selected.FullPath() != "" {
						secretInfo, err := s.gcp.GetSecretInfo(selected.FullPath())
						if err != nil {
							log.Error().Err(err).Msg("Error getting secret info")
							s.components.toast.SetText("Error getting secret info")
						} else {
							s.Modal = view.NewSecretInfoModal(secretInfo, selected)
							s.Modal.Init()
						}
					} else {
						s.components.toast.SetText("No secret selected")
					}
				case "tab":
					s.components.list.ToggleFocus()
					s.components.detail.ToggleFocus()
					s.components.detail.SetFilteredValue(s.components.list.SearchQuery)
				case "shift+right":
					cmds = append(cmds, adjustListWidth(s, 1))
					return tea.Batch(cmds...)
				case "shift+left":
					cmds = append(cmds, adjustListWidth(s, -1))
					return tea.Batch(cmds...)
				}
			}
		case editor.EditFinishedMsg:
			if msg.Equal {
				s.components.toast.SetText("No changes detected")
			} else {
				log.Info().Msgf("Changes detected in secret %v", msg.CurrentSecret.Title())
				s.components.toast.SetText("Changes detected")

				s.Modal = view.NewConfirm("Do you want to create a new secret based on this?", msg)
				s.Modal.Init()

				s.components.detail.SetContent(string(msg.SecretData))
				_, cmd = s.components.detail.Update(msg)
				return cmd
			}
			return nil
		case SecretLoadedMsg:
			s.components.detail.SetContent(msg.Text)
			return nil
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				s.Modal = nil
				return nil
			}
		}

		modal, cmd := s.Modal.Update(msg)
		s.Modal = modal
		return cmd
	}

	_, cmd = s.components.list.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = s.components.detail.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = s.components.help.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = s.components.toast.Update(msg)
	cmds = append(cmds, cmd)

	if s.components.list.IsFiltering() == false {
		cmd = s.showSecret()
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

type SecretLoadedMsg struct {
	Secret view.Secret
	Text   string
}

func adjustListWidth(s *Secrets, delta int) tea.Cmd {
	s.ListWidth += delta
	if s.ListWidth > 50 {
		s.ListWidth = 50
	} else if s.ListWidth < 6 {
		s.ListWidth = 6
	}
	resizeCmd := func() tea.Msg {
		return view.ResizeMessage{}
	}
	return resizeCmd
}

func (s *Secrets) showSecret() tea.Cmd {
	selected := s.components.list.SelectedItem()

	if selected.FullPath() == "" {
		return func() tea.Msg {
			return SecretLoadedMsg{
				Secret: selected,
				Text:   "",
			}
		}
	}

	return func() tea.Msg {
		var text string
		text = "loading"
		if selected.Type() == "version" {
			versionSecret, err := s.gcp.GetSecretVersion(selected.FullPath(), strconv.Itoa(selected.Version()))
			if err != nil {
				text = "Error loading secret version: " + err.Error()
			} else {
				text = ui.SyntaxHighlight(versionSecret)
			}
		} else {
			secretData, err := s.gcp.GetSecret(selected.FullPath())
			if err != nil {
				text = "Error loading secret: " + err.Error()
			} else {
				text = ui.SyntaxHighlight(secretData)
			}
		}
		return SecretLoadedMsg{
			Secret: selected,
			Text:   text,
		}
	}

}

func NewSecrets(gcp client.Client, selected int) *Secrets {
	page := &Secrets{gcp: gcp, ListWidth: 31}
	page.Init()
	page.Select(selected)
	return page
}

func (s *Secrets) Init() {
	secretList := view.NewSecretsList(50, 50, s.gcp)
	secretView := view.NewSecretView(50, 50)
	help := view.NewHelp()
	toast := view.NewToast()

	s.components = secretsComponents{
		list:   &secretList,
		detail: &secretView,
		help:   &help,
		toast:  &toast,
	}

	s.Update(s.showSecret()())
}

func (s *Secrets) Select(index int) {
	s.components.list.Select(index)
}
