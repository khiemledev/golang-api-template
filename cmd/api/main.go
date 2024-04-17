package main

import (
	"khiemle.dev/golang-api-template/api/routes"
	util "khiemle.dev/golang-api-template/pkg/util"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load config")
		return
	}

	util.ConfigLogger(cfg)

	r := routes.SetupRouter()
	err = r.Run(cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("Error while starting server")
	}

	log.Info().Msgf("Listening and serving HTTP on %s", cfg.HTTPServerAddress)
}
