package main

import (
	"khiemle.dev/golang-api-template/api"
	"khiemle.dev/golang-api-template/pkg/database"
	util "khiemle.dev/golang-api-template/pkg/util"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load config")
		return
	}

	// Initialize a new GORM DB instance
	db, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error initializing GORM DB: %v\n", err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting SQL DB")
		return
	}
	defer sqlDB.Close()
	log.Info().Msg("Connected to database")

	// Run database migrations
	err = database.DBMigration(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not run database migrations")
		return
	}
	log.Info().Msg("Ran database migrations")

	// Create API server
	server := api.NewServer()
	err = server.Initialize(&cfg, db)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not initialize server")
		return
	}
	log.Info().Msg("Initialized server")

	// Start server
	err = server.StartServer()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start server")
		return
	}
	log.Info().Msgf("Listening and serving HTTP on %s", cfg.HTTPServerAddress)
}
