package main

import (
	"github.com/gin-gonic/gin"
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
	//service := jcService.Init(config, logger, axiosClient)

	r := gin.Default()
	r.POST("/cage", service.CreateCage)
	r.GET("/cage", service.GetAllCages)
	r.GET("/cage/:id", service.GetCage)
	r.GET("/dinosaur/:id", service.GetDinosaur)
	r.POST("dinosaur", service.CreateDinosaur)
	r.GET("/dinosaur", service.GetAllDinosaurs)
	r.PUT("/dinosaur/:id/:cage_id", service.AddDinosaurToCage)
	r.Run(config.HTTPListen)
}
