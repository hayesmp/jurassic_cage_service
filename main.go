package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hayesmp/jurassic-cage-service/internal"
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

	_, err = internal.ParseConfigFromEnvironment("")
	if err != nil {
		logger.Fatal().Err(err).Msg("failure parsing config")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
