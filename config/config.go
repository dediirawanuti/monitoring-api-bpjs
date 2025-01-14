package main

import (
	"os"
)

// Configurations for the API and Telegram Bot
type Config struct {
	APIUrl       string
	TelegramBot  string
	TelegramChat string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() Config {
	return Config{
		APIUrl:       os.Getenv("API_URL"),
		TelegramBot:  os.Getenv("TELEGRAM_BOT"),
		TelegramChat: os.Getenv("TELEGRAM_CHAT"),
	}
}
