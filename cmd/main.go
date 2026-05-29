package main

import (
	"flag"
	"fmt"
	"os"
	"smm/internal/bootstrap"
	"smm/internal/config"
	"smm/internal/model"

	tea "charm.land/bubbletea/v2"
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

	p := tea.NewProgram(model.New(projectId))

	_, err := p.Run()
	if err != nil {
		return
	}
}
