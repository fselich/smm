package main

import (
	"flag"
	"fmt"
	"os"
	"smm/internal/bootstrap"
	"smm/internal/config"
	"smm/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

var version = "dev"

func init() {
	bootstrap.LoadConfig()
	bootstrap.SetLog()
}

func main() {
	projectIdFlag := flag.String("p", "", "Project ID to use")
	versionFlag := flag.Bool("v", false, "Show version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

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
