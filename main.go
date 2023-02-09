package main

import (
	"github.com/hayesmp/jurassic-cage-service/internal"
	jcService "github.com/hayesmp/jurassic-cage-service/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"
)

var logger = zerolog.New(os.Stderr)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	config, err := internal.ParseConfigFromEnvironment("")
	if err != nil {
		logger.Fatal().Err(err).Msg("failure parsing config")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	service := jcService.Init(config, logger)

	r := service.SetupRouter()
	r.Run(config.HTTPListen)
}
