package main

import (
	"gcs/bootstrap"
	"gcs/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

func init() {
	bootstrap.SetLog()
	bootstrap.LoadConfig()
}

func main() {
	projectId := viper.GetString("selected")
	p := tea.NewProgram(model.New(projectId), tea.WithAltScreen())

	_, err := p.Run()
	if err != nil {
		return
	}
}
