package bootstrap

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"path"
)

func LoadEnv() {
	envFiles := []string{
		".env",
		fmt.Sprintf(".env.%s", os.Getenv("RUN_ENV")),
		".env.local",
	}

	envLoaded := false
	for _, envFile := range envFiles {
		if err := godotenv.Overload(ConfigPath(envFile)); err != nil {
			log.Debug().Err(err).Msg(fmt.Sprintf("Error loading %s file", envFile))
		} else {
			log.Info().Msg(fmt.Sprintf("Loaded %s file", ConfigPath(envFile)))
			envLoaded = true
		}
	}

	if !envLoaded {
		log.Fatal().Msg("No env file loaded")
	}
}

func ConfigPath(file string) string {
	filename, _ := os.Executable()
	return path.Join(path.Dir(filename), "/config", file)
}

func SetLog() {
	file, err := os.OpenFile("myapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}

	log.Logger = log.Output(file)
	log.Info().Msg("Log file opened")
}
