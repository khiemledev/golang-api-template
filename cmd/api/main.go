package main

import (
	"khiemle.dev/golang-api-template/api"
	"khiemle.dev/golang-api-template/pkg/database"
	util "khiemle.dev/golang-api-template/pkg/util"

	"github.com/rs/zerolog/log"

	_ "khiemle.dev/golang-api-template/docs"
)

//	@title			Golang API Template
//	@version		0.1.0
//	@description	This is the template for Golang API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Khiem Le
//	@contact.url	https://khiemle.dev
//	@contact.email	khiemledev@gmail.com

//	@host		localhost:8085
//	@BasePath	/v1

//	@securityDefinitions.apiKey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter JWT Bearer token
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
