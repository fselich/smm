package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Project struct {
	ID   string `yaml:"id" json:"id"`
	Type string `yaml:"type" json:"type"`
}

func Load() error {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "smm")
	configFile := filepath.Join(configPath, "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, 0700)
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating config directory")
			return err
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetDefault("projects", []Project{})

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err = viper.WriteConfigAs(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing config file")
			return err
		}
	} else {
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Error reading config file")
			return err
		}
	}

	return nil
}

func Save() error {
	return viper.WriteConfig()
}

func GetSelectedProjectId() string {
	return viper.GetString("selected")
}

func SetSelectedProject(projectId string) {
	viper.Set("selected", projectId)
}

func GetProjectIDs() []string {
	var projects []Project
	err := viper.UnmarshalKey("projects", &projects)
	if err != nil {
		return []string{}
	}

	ids := make([]string, len(projects))
	for i, project := range projects {
		ids[i] = project.ID
	}
	return ids
}

func AddProjectID(projectId string) {
	var projects []Project
	err := viper.UnmarshalKey("projects", &projects)
	if err != nil {
		projects = []Project{}
	}

	for _, project := range projects {
		if project.ID == projectId {
			return
		}
	}

	projects = append(projects, Project{ID: projectId, Type: "gcp"})
	viper.Set("projects", projects)
}

func GetTypeByProjectId(projectId string) string {
	var projects []Project
	err := viper.UnmarshalKey("projects", &projects)
	if err != nil {
		return ""
	}

	for _, project := range projects {
		if project.ID == projectId {
			return project.Type
		}
	}
	return "gcp"
}

func GetLogPath() string {
	return viper.GetString("logPath")
}
