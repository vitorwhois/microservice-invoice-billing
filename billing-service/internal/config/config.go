package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Database            DatabaseConfig
	Server              ServerConfig
	InventoryServiceURL string
	DatabaseURL         string
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
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env não encontrado, usando variáveis do ambiente")
	}
	viper.AutomaticEnv()
	viper.SetDefault("DB_HOST", "billing-db")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "billing")
	viper.SetDefault("SERVER_PORT", "8081")
	viper.SetDefault("INVENTORY_SERVICE_URL", "http://inventory-service:8080")

	return &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
	}, nil
}
