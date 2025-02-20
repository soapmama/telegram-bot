package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Token string
	Links struct {
		Distillate string `mapstructure:"distillate"`
		Prices     string `mapstructure:"prices"`
		Soap       string `mapstructure:"soap"`
	} `mapstructure:"links"`
}

func newConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	botToken := v.GetString("TOKEN")
	if botToken == "" {
		log.Fatal("TOKEN environment variable not set")
	}

	config := &Config{}
	if err := v.Unmarshal(config); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}

	if config.Token == "" {
		log.Fatal("TOKEN environment variable not set")
	}

	return config
}
