package config

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	MessageBackend  string `json:"message_backend"`
	BackendSettings struct {
		Telegram struct {
			Token  string `json:"token"`
			ChatID string `json:"chat_id"`
		} `json:"telegram"`
		Gotify struct {
			Token string `json:"token"`
			Url   string `json:"url"`
		}
	} `json:"backend_settings"`
}

func DecodeConfig(configPath string) (*Config, error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Warn("No config file found, creating one")
		err = setupConfig(configPath)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	defer configFile.Close()

	config, err := readConfig(configFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func readConfig(f *os.File) (decodedConfig *Config, err error) {
	decoder := json.NewDecoder(f)
	decodedConfig = &Config{}
	err = decoder.Decode(&decodedConfig)

	return
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
