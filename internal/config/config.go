package config

import (
	"encoding/json"
	"os"
)

func Decode(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	decodedConfig := &Config{}
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return nil, err
	}
	return decodedConfig, nil
}

func validateConfig(config *Config) error {
	return nil
	// check that file exists
	// if not exist, create it
	// check that all required values are present
	// if not, enter interactive setup
}

func setupConfig(configPath string) error {
	return nil
	// ask which backend to use
	// ask for required values depending on selected backend
	// write data to file
}

type Config struct {
	MessageBackends MessageBackends `json:"message_backends"`
}

type MessageBackends struct {
	Telegram MessageBackendTelegram `json:"telegram"`
}

type MessageBackendTelegram struct {
	Enabled         bool   `json:"enabled"`
	Token           string `json:"token"`
	ChatID          string `json:"chat_id"`
	MessageThreadID string `json:"message_thread_id"`
}
