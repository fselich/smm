package page

import (
	gcp2 "gcs/gcp"
	"gcs/ui"
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

type Versions struct {
	secret     secret
	title      string
	gcp        *gcp2.Gcp
	index      int
	components map[string]any
}

type secret interface {
	Name() string
	Title() string
	Index() int
}

func NewVersions(secret secret, gcp *gcp2.Gcp) *Versions {
	page := &Versions{secret: secret, gcp: gcp}
	page.Init()
	return page
}

func (v *Versions) Init() {
	secretList := view.NewSecretVersionsList(50, 50, v.secret.Name(), v.secret.Title(), v.gcp)
	secretView := view.NewSecretView(50, 50)
	help := view.NewHelp()
	toast := view.NewToast()
	v.components = make(map[string]any)
	v.components["secretList"] = &secretList
	v.components["secretView"] = &secretView
	v.components["help"] = &help
	v.components["toast"] = &toast
}

func (v *Versions) View() string {
	list := v.components["secretList"].(*view.SecretVersionsList)
	detail := v.components["secretView"].(*view.SecretView)
	help := v.components["help"].(*view.Help)
	toast := v.components["toast"].(*view.Toast)
	listView := list.View()
	detailView := detail.View()
	borderColor := lipgloss.Color("#87CEFA")
	borderedList := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Width(list.Width()).
		Render(listView)
	borderedDetail := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor).
		Render(detailView)

	borderedHelp := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#C0C0C0")).
		Width(list.Width() + detail.Width() + 2).
		Render(help.View())

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Top, borderedList, borderedDetail),
		lipgloss.JoinHorizontal(lipgloss.Bottom, toast.View()),
		lipgloss.JoinHorizontal(lipgloss.Bottom, borderedHelp))
}

func (v *Versions) Resize(width int, height int) {
	log.Info().Msgf("[status 3] dimensions: %v, %v", width, height)
	list := v.components["secretList"].(*view.SecretVersionsList)
	detail := v.components["secretView"].(*view.SecretView)
	help := v.components["help"].(*view.Help)

	list.SetHeight(height - 9)
	list.SetWidth(30)
	detail.SetWidth(width - 4 - list.Width())
	detail.SetHeight(height - 9)
	help.SetWidth(list.Width() + detail.Width() + 2)
}

func (v *Versions) Update(msg tea.Msg) tea.Cmd {
	list := v.components["secretList"].(*view.SecretVersionsList)
	detail := v.components["secretView"].(*view.SecretView)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			cmd = func() tea.Msg {
				return SetStatusMsg{Status: 2, From: "versions", Data: v.secret}
			}
			return cmd
		}
	}
	*list, cmd = list.Update(msg)
	cmds = append(cmds, cmd)

	selected := list.SelectedItem()
	secretData := v.gcp.GetSecretVersion(selected.FullPath(), selected.Description())
	text := ui.SyntaxHighlight(secretData)
	detail.SetContent(text)
	*detail, cmd = detail.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}
