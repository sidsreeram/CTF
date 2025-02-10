package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DATABASEURL string `mapstructure:"DATABASE_URL"`
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("../../cmd/server/.env")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	// Ensure that the DBURL field is also set
	if config.DATABASEURL == "" {
		return config, fmt.Errorf("DB_URL is required in the configuration")
	}
	return config, nil
}
