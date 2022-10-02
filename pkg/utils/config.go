package utils

import (
	"github.com/spf13/viper"
)

// Config store all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	PORT      string `mapstructure:"PORT"`
	MONGO_URI string `mapstructure:"MONGO_URI"`
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
