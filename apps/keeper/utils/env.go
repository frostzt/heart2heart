package utils

import (
	"os"

	"github.com/spf13/viper"
)

// Env has environment stored
type Env struct {
	Environment string `mapstructure:"ENV"`
	ServerPort  string `mapstructure:"SERVER_PORT"`

	// Database envs
	DB_HOST string `mapstructure:"DB_HOST"`
	DB_PORT string `mapstructure:"DB_PORT"`
	DB_NAME string `mapstructure:"DB_NAME"`
	DB_USER string `mapstructure:"DB_USER"`
	DB_PASS string `mapstructure:"DB_PASS"`

	// Authentication
	JWT_SECRET        string `mapstructure:"JWT_SECRET"`
	JWT_COOKIE_DOMAIN string `mapstructure:"JWT_COOKIE_DOMAIN"`

	// OpenTelemetry and Signoz
	SERVICE_NAME                string `mapstructure:"SERVICE_NAME"`
	OTEL_EXPORTER_OTLP_ENDPOINT string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	INSECURE_MODE               string `mapstructure:"INSECURE_MODE"`
}

// NewEnv creates a new environment
func NewEnv(log Logger) Env {
	env := Env{
		Environment:                 os.Getenv("ENV"),
		ServerPort:                  os.Getenv("SERVER_PORT"),
		DB_HOST:                     os.Getenv("DB_HOST"),
		DB_PORT:                     os.Getenv("DB_PORT"),
		DB_NAME:                     os.Getenv("DB_NAME"),
		DB_USER:                     os.Getenv("DB_USER"),
		DB_PASS:                     os.Getenv("DB_PASS"),
		JWT_SECRET:                  os.Getenv("JWT_SECRET"),
		JWT_COOKIE_DOMAIN:           os.Getenv("JWT_COOKIE_DOMAIN"),
		SERVICE_NAME:                os.Getenv("SERVICE_NAME"),
		OTEL_EXPORTER_OTLP_ENDPOINT: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		INSECURE_MODE:               os.Getenv("INSECURE_MODE"),
	}
	viper.AutomaticEnv()

	err := viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	return env
}
