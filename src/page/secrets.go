package page

import (
	"fmt"
	"gcs/editor"
	gcp2 "gcs/gcp"
	"gcs/ui"
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Page interface {
	Init()
	View() string
	Resize(int, int)
	Update(cmd tea.Msg) tea.Cmd
}

type SetStatusMsg struct {
	Status int
	From   string
	Data   any
}

type Secrets struct {
	gcp        *gcp2.Gcp
	components map[string]any
	Modal      view.Modal
}

type CurrentSecret struct {
	name  string
	title string
	index int
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
	list := s.components["list"].(*view.SecretsList)
	detail := s.components["detail"].(*view.SecretView)
	help := s.components["help"].(*view.Help)
	toast := s.components["toast"].(*view.Toast)

	listView := list.View()
	detailView := detail.View()

	borderedList := ui.StyleBorder().
		Width(list.Width()).
		Render(listView)
	borderedDetail := ui.StyleBorder().
		Render(detailView)

	borderedHelp := ui.StyleLowBorder().
		Width(list.Width() + detail.Width() + 2).
		Render(help.View())

	var x int
	if list.IsFiltering() {
		x = ((list.Width() - len(list.FilterValue())) / 2) - 1
		listTitle := ui.StyleBorderTitle().Render(list.FilterValue())
		borderedList = ui.PlaceOverlay(x, 0, "\""+listTitle+"█\"", borderedList, false)
	} else if list.IsFiltered() {
		x = (list.Width() - len(list.FilterValue())) / 2
		listTitle := ui.StyleBorderTitle().Render(list.FilterValue())
		borderedList = ui.PlaceOverlay(x, 0, listTitle, borderedList, false)
	}

	x = (detail.Width() - len(list.SelectedItem().Title())) / 2
	detailTitle := ui.StyleBorderTitle().Render(list.SelectedItem().Title())
	borderedDetail = ui.PlaceOverlay(x, 0, detailTitle, borderedDetail, false)

	render := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, borderedList, borderedDetail),
		lipgloss.JoinHorizontal(lipgloss.Bottom, borderedHelp),
		lipgloss.JoinHorizontal(lipgloss.Bottom, toast.View()),
	)

	if s.Modal != nil {
		render = ui.ModalOverlay(s.Modal.View(), render)
	}

	return render
}

func (s *Secrets) Resize(width int, height int) {
	list := s.components["list"].(*view.SecretsList)
	detail := s.components["detail"].(*view.SecretView)
	help := s.components["help"].(*view.Help)
	toast := s.components["toast"].(*view.Toast)

	if list.IsFiltering() {
		list.SetHeight(height - 12)
	} else {
		list.SetHeight(height - 6)
	}

	list.SetWidth(30)
	detail.SetWidth(width - 4 - list.Width())
	detail.SetHeight(height - 6)
	help.SetWidth(list.Width() + detail.Width() + 2)
	toast.SetWith(width)
}

func (s *Secrets) Update(msg tea.Msg) tea.Cmd {
	list := s.components["list"].(*view.SecretsList)
	detail := s.components["detail"].(*view.SecretView)
	help := s.components["help"].(*view.Help)
	toast := s.components["toast"].(*view.Toast)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case view.ProjectSelectedMessage:
		log.Info().Msgf("Project selected: %v", msg.ProjectId)
		s.Modal = nil
		return nil
	}

	if s.Modal == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if list.IsFiltering() == false {
				switch msg.String() {
				case "e":
					f, err := os.Create("detail_content.env")
					if err != nil {
						log.Fatal()
					}
					defer f.Close()

					secretName := list.SelectedItem().FullPath()
					secretData := string(s.gcp.GetSecret(secretName))
					_, err = f.WriteString(secretData)
					if err != nil {
						log.Fatal()
					}

					return editor.OpenEditor(secretData, list.SelectedItem())
				case "p":
					s.Modal = view.NewProjectSelectorModal()
					s.Modal.Init()
				case "ctrl+c":
					return tea.Quit
				case "n":
					err := s.gcp.AddSecretVersion(list.SelectedItem().Title(), []byte("prueba2"))
					if err != nil {
						log.Error().Msgf("Error creating new secret version: %v", err)
					}
				case "ctrl+n":
					log.Info().Msg("New secret version")

				case "y":
					if s.Modal == nil {
						s.Modal = view.NewConfirm("Do you want to create a new secret based on this?", "hola")
						s.Modal.Init()
					} else {
						s.Modal = nil
					}
				case "h":
					detail.Hidden = !detail.Hidden
				case "v":
					selected := list.SelectedItem()
					if selected.Type() == "current" {
						deleted := list.DelVersionItems()
						if !deleted {
							versions := s.gcp.GetSecretVersions(selected.FullPath())
							versions = versions[1:]
							toast.SetText(fmt.Sprintf("Secret has %v versions", len(versions)))
							for i, version := range versions {
								secret := view.NewSecret(strconv.Itoa(version.Version), version.FullPath, "version", version.Version)
								secret.SetRelated(&selected)
								cmd = list.InsertItem(list.RealIndex()+1+i, secret)
							}
						}
					}
					list.Select(selected.Index())
					return cmd
				case "c":
					toast.SetText("Secret copied to clipboard")
				case "r":
					s.Init()
					toast.SetText("Secrets refreshed")
				}
			}
		case editor.EditorFinishedMsg:
			toast.SetText("Saved")
			if msg.Err != nil {
				return tea.Quit
			}
			return nil
		case SecretLoadedMsg:
			detail.SetContent(msg.Text)
			return nil
		}
	} else {
		modal, cmd := s.Modal.Update(msg)
		s.Modal = modal
		return cmd
	}

	_, cmd = list.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = detail.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = help.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = toast.Update(msg)
	cmds = append(cmds, cmd)

	if list.IsFiltering() == false {
		cmd = s.showSecret()
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

type SecretLoadedMsg struct {
	Secret view.Secret
	Text   string
}

func (s *Secrets) showSecret() tea.Cmd {
	list := s.components["list"].(*view.SecretsList)
	detail := s.components["detail"].(*view.SecretView)

	if detail.Hidden {
		detail.SetContent("")
		return nil
	}

	selected := list.SelectedItem()

	return func() tea.Msg {
		var text string
		text = "loading"
		if selected.Type() == "version" {
			versionSecret := s.gcp.GetSecretVersion(selected.FullPath(), strconv.Itoa(selected.Version()))
			text = ui.SyntaxHighlight(versionSecret)
		} else {

			secretData := s.gcp.GetSecret(selected.FullPath())
			text = ui.SyntaxHighlight(secretData)
		}
		return SecretLoadedMsg{
			Secret: selected,
			Text:   text,
		}
	}

}

func NewSecrets(gcp *gcp2.Gcp, selected int) *Secrets {
	page := &Secrets{gcp: gcp}
	page.Init()
	page.Select(selected)
	return page
}

func (s *Secrets) Init() {
	secretList := view.NewSecretsList(50, 50, s.gcp)
	secretView := view.NewSecretView(50, 50)
	help := view.NewHelp()
	toast := view.NewToast()
	s.components = make(map[string]any)
	s.components["list"] = &secretList
	s.components["detail"] = &secretView
	s.components["help"] = &help
	s.components["toast"] = &toast
	s.showSecret()
}

func (s *Secrets) Select(index int) {
	s.components["list"].(*view.SecretsList).Select(index)
}
