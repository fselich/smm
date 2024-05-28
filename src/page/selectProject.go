package page

import (
	"fmt"
	"gcs/view"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog/log"
)

type component interface {
	View() string
	Update(msg tea.Msg) (component, tea.Cmd)
}

type Page interface {
	Init()
	View() string
	Resize(int, int)
	Update(cmd tea.Msg) tea.Cmd
}

type ProjectSelectedMsg struct {
	ProjectId string
}

type SetStatusMsg struct {
	Status int
	From   string
	Data   any
}

type SelectProject struct {
	components map[string]any
}

func (s *SelectProject) View() string {
	selector := s.components["selector"].(*view.ProjectSelector)
	return fmt.Sprintf(
		"Project ID?\n\n%s\n\n%s",
		selector.View(),
		"(esc to quit)",
	) + "\n"
}

func (s *SelectProject) Resize(width int, height int) {
	//TODO implement me
	log.Info().Msg("Resize")
}

func (s *SelectProject) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	selector := s.components["selector"].(*view.ProjectSelector)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd = func() tea.Msg {
				return ProjectSelectedMsg{ProjectId: selector.Value()}
			}
		default:
			_, cmd = selector.Update(msg)
		}
	}

	return cmd
}

func NewSelectProject() *SelectProject {
	page := SelectProject{}
	page.Init()
	return &page
}

func (s *SelectProject) Init() {
	s.components = make(map[string]any)
	selector := view.NewProjectSelector()
	s.components["selector"] = &selector
}
