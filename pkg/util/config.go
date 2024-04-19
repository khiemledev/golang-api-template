package util

import (
	"github.com/spf13/viper"
)

// Config stores all application configs loaded from file or env variables
type Config struct {
	// HTTP Server
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`

	// Swagger
	SwaggerURL string `mapstructure:"SWAGGER_URL"`

	// Logging
	LogFilename   string `mapstructure:"LOG_FILENAME"`
	LogMaxSize    int    `mapstructure:"LOG_MAX_SIZE"` // in megabytes
	LogMaxBackups int    `mapstructure:"LOG_MAX_BACKUPS"`
	LogMaxAge     int    `mapstructure:"LOG_MAX_AGE"`
	LogCompress   bool   `mapstructure:"LOG_COMPRESS"`

	// Database
	DBHost     string `mapstructure:"DB_HOST"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     int    `mapstructure:"DB_PORT"`

	// JWT
	PasetoSecretHex string `mapstructure:"JWT_SECRET_KEY"`
}

// Load config from env file and return Config
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("HTTP_SERVER_ADDRESS", ":8080")

	viper.SetDefault("SWAGGER_URL", "/docs")

	// Configs for logger
	viper.SetDefault("LOG_FILENAME", "logs/app.log")
	viper.SetDefault("LOG_MAX_SIZE", 10)
	viper.SetDefault("LOG_MAX_BACKUPS", 5)
	viper.SetDefault("LOG_MAX_AGE", 28)
	viper.SetDefault("LOG_COMPRESS", true)

	// Read config
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal config
	err = viper.Unmarshal(&config)
	return
}
