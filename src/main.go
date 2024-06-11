package main

import (
	"gcs/bootstrap"
	"gcs/model"
	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	bootstrap.SetLog()
	bootstrap.LoadConfig()
}

func main() {
	m := model.New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return
	}
}
