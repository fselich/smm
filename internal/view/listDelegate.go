package view

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/truncate"
	"io"
	"smm/internal/ui"
)

type ItemDelegate struct {
	Styles list.DefaultItemStyles
}

func NewListDelegate() *ItemDelegate {
	return &ItemDelegate{
		Styles: list.NewDefaultItemStyles(),
	}
}

const (
	ellipsis = "…"
)

func (d *ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		s = &d.Styles
	)

	title := item.(Secret).Title()
	if item.(Secret).Type() == "version" {
		if item.(Secret).Title() == "1" {
			title = fmt.Sprintf("%s%s [v.%s]", ui.StyleLow().Render("└──"), item.(Secret).CreatedAt().Format("2006-01-02 15:04:05"), title)
		} else {
			title = fmt.Sprintf("%s%s [v.%s]", ui.StyleLow().Render("├──"), item.(Secret).CreatedAt().Format("2006-01-02 15:04:05"), title)
		}

	}
	textWidth := uint(m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight())
	title = truncate.StringWithTail(title, textWidth, ellipsis)

	isSelected := index == m.Index()

	if isSelected {
		title = s.SelectedTitle.Padding(0, 0, 0, 0).Render(title)
	} else {
		title = s.NormalTitle.Padding(0, 0, 0, 1).Render(title)
	}

	_, _ = fmt.Fprintf(w, "%s", title)
}

func (d *ItemDelegate) Height() int {
	return 1
}

func (d *ItemDelegate) Spacing() int {
	return 0
}

func (d *ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
