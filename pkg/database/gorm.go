package database

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_authModel "khiemle.dev/golang-api-template/internal/auth/model"
	_todoModel "khiemle.dev/golang-api-template/internal/todo/model"
	_userModel "khiemle.dev/golang-api-template/internal/user/model"
	util "khiemle.dev/golang-api-template/pkg/util"
)

// NewGormDB initializes a new GORM DB instance with the provided database URL.
func NewGormDB(cfg util.Config) (*gorm.DB, error) {
	// Connect to the PostgreSQL database
	psqlconn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	) // noqa
	db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting SQL DB")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

func DBMigration(db *gorm.DB) error {
	var err error
	err = db.AutoMigrate(&_todoModel.Todo{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&_userModel.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&_authModel.LoginSession{})
	if err != nil {
		return err
	}

	return nil
}
