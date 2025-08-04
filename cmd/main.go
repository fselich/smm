package main

import (
	"flag"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"smm/internal/bootstrap"
	"smm/internal/model"
)

func init() {
	bootstrap.SetLog()
	bootstrap.LoadConfig()
}

func main() {
	projectIdFlag := flag.String("p", "", "Project ID to use")
	flag.Parse()

	projectId := *projectIdFlag
	if projectId == "" {
		projectId = viper.GetString("selected")
	}

	p := tea.NewProgram(model.New(projectId), tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return
	}
}
