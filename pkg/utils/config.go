package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// Config store all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	PORT      string `mapstructure:"PORT"`
	MONGO_URI string `mapstructure:"MONGO_URI"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig() (config Config, err error) {
	err = godotenv.Load(".env")
	config.PORT = os.Getenv("PORT")
	config.MONGO_URI = os.Getenv("MONGO_URI")

	return
}
