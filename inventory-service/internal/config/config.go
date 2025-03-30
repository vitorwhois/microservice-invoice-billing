package config

import "github.com/spf13/viper"

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
}

func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetDefault("DB_HOST", "localhost")

	return &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Server: ServerConfig{
			Port:         viper.GetString("PORT"),
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
	}, nil
}
