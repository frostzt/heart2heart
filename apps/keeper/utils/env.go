package utils

import (
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
}

// NewEnv creates a new environment
func NewEnv(log Logger) Env {
	env := Env{}
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("☠️ cannot read configuration")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("☠️ environment can't be loaded: ", err)
	}

	log.Infof("%+v \n", env)
	return env
}
