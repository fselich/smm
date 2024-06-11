package bootstrap

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func LoadConfig() {
	configPath := filepath.Join(os.Getenv("HOME"), ".config", "gcs")
	configFile := filepath.Join(configPath, "config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("Error creating config directory")
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.SetDefault("projectIds", []string{"civitatis-guias-1", "otra cosa"})

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err = viper.WriteConfigAs(configFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Error writing config file")
		}
	} else {
		err = viper.ReadInConfig()
		if err != nil {
			log.Fatal().Err(err).Msg("Error reading config file")
		}
	}
}

func SetLog() {
	file, err := os.OpenFile("gcs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	log.Logger = log.Output(file)
	log.Info().Msg("Log file opened")
}
