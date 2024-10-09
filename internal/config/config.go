package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database   DatabaseConfig
	Processing ProcessingConfig
	Metrics    MetricsConfig
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type ProcessingConfig struct {
	BatchSize          int
	MaxConcurrentFiles int
}

type MetricsConfig struct {
	Port int
}

func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
