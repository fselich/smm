package main

import (
	"flag"
	"smm/internal/bootstrap"
	"smm/internal/config"
	"smm/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	bootstrap.LoadConfig()
	bootstrap.SetLog()
}

func main() {
	projectIdFlag := flag.String("p", "", "Project ID to use")
	flag.Parse()

	projectId := *projectIdFlag
	if projectId == "" {
		projectId = config.GetSelectedProjectId()
	}

	p := tea.NewProgram(model.New(projectId), tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return
	}
}
