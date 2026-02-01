package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	App    AppConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int
	Host string
}

// AppConfig holds application-related configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

var AppConfiguration *Config

// LoadConfig loads configuration from .env file and environment variables
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_HOST", "localhost")
	viper.SetDefault("APP_NAME", "Kasir API")
	viper.SetDefault("APP_VERSION", "1.0")
	viper.SetDefault("APP_ENVIRONMENT", "development")

	// Read environment variables
	viper.AutomaticEnv()

	// Read .env file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(".env file not found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading .env file: %w", err)
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("SERVER_PORT"),
			Host: viper.GetString("SERVER_HOST"),
		},
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Version:     viper.GetString("APP_VERSION"),
			Environment: viper.GetString("APP_ENVIRONMENT"),
		},
	}

	AppConfiguration = config
	return config, nil
}

// GetServerAddress returns the server address in host:port format
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
