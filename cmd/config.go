package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Links struct {
	Distillate string `mapstructure:"distillate"`
	Prices     string `mapstructure:"prices"`
	Soap       string `mapstructure:"soap"`
	Ubtan      string `mapstructure:"ubtan"`
}

type Config struct {
	Token    string `mapstructure:"TOKEN"`
	Port     string `mapstructure:"PORT"`
	ChatID   int64  `mapstructure:"CHAT_ID"`
	ThreadID int64  `mapstructure:"THREAD_ID"`
	Links    Links  `mapstructure:"links"`
}

func newConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file %v", err)
	}

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.BindEnv("TOKEN")
	v.BindEnv("CHAT_ID")
	v.BindEnv("THREAD_ID")
	v.BindEnv("PORT")

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

	log.Printf("Config: %+v", v.AllKeys())
	if config.Token == "" {
		log.Fatal("TOKEN environment variable not set")
	}

	return config
}
