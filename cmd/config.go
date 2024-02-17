package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ContextKey string

type Config struct {
	AppEnv string `mapstructure:"APPENV"`
	AppTz  string `mapstructure:"TZ"`
	// AppIsDev          bool
	DATABASE_URL      string `mapstructure:"DATABASE_URL"`
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_DB       string `mapstructure:"POSTGRES_DB"`
}

func NewConfig() (*Config, error) {
	env := os.Getenv("APPENV")
	if env == "" {
		env = "local"
	}

	viper.AddConfigPath(".")
	viper.AddConfigPath("cmd")
	viper.SetConfigName(env)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetConfigName("local")

			if err := viper.ReadInConfig(); err != nil {
				log.Error().Err(err).Msg("Failed To Read Config")
				return nil, err
			}
		} else {
			log.Error().Err(err).Msg("Failed To Read Config")
			return nil, err
		}
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Error().Err(err).Msg("Failed To Unmarshal Config")
		return nil, err
	}

	// cfg.AppIsDev = cfg.AppEnv == "staging" || cfg.AppEnv == "local" || cfg.AppEnv == "dev"

	return cfg, nil
}
