package main

import (
	"khiemle.dev/golang-api-template/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load config")
		return
	}

	utils.ConfigLogger(cfg)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		log.Info().Msg("Ping request received!")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err = r.Run(cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("Error while starting server")
	}

	log.Info().Msgf("Listening and serving HTTP on %s", cfg.HTTPServerAddress)
}
