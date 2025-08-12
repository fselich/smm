package bootstrap

import (
	"os"
	"smm/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoadConfig() {
	config.Load()
}

func SetLog() {
	logPath := config.GetLogPath()

	if logPath == "" {
		log.Logger = log.Logger.Level(zerolog.Disabled)
		return
	}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening log file")
	}

	log.Logger = log.Output(file)
}
