package util

import (
	"os"

	"github.com/joho/godotenv"
)

// Config store all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	Port      string `mapstructure:"PORT"`
	MongoUri  string `mapstructure:"MONGO_URI"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

var AppConfig Config

// LoadConfig reads configuration from file or environment variables
func LoadConfig() {
	_ = godotenv.Load(".env")
	AppConfig.Port = os.Getenv("PORT")
	AppConfig.MongoUri = os.Getenv("MONGO_URI")
	AppConfig.SecretKey = os.Getenv("SECRET_KEY")
}
