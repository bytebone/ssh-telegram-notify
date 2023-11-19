package tgsend

// read config file
// determine host os (win/posix)
// find config file
// read config file
// parse bot token, group id
// send telegram message

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	BotToken string
	GroupID  string
}

var configPath string = getHomeDir() + "/.config/tgsend/config.json"

func SendMessage(message string) error {
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Warn("No config file found, creating one")
		err = CreateEmptyConfig()
		if err != nil {
			return err
		}
		return err
	}
	defer configFile.Close()

	config, err := readConfig(configFile)
	if err != nil {
		return err
	}

	switch {
	case config.BotToken == "":
		return errors.New("no telegram bot token found")
	case config.GroupID == "":
		return errors.New("no telegram group id found")
	case config.BotToken == "" && config.GroupID == "":
		return errors.New("no telegram bot token or group id found")
	}

	data := url.Values{
		"chat_id":    {config.GroupID},
		"text":       {message},
		"parse_mode": {"Markdown"},
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)
	res, err := http.PostForm(url, data)
	log.Debug(url)
	log.Debug(data)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", res.Status)
	}
	return err
}

func readConfig(f *os.File) (decodedConfig Config, err error) {
	decoder := json.NewDecoder(f)
	decodedConfig = Config{}
	err = decoder.Decode(&decodedConfig)
	return
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return homeDir
}

func CreateEmptyConfig() error {
	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	_, err = configFile.Write([]byte(`{"BotToken":"","GroupID":""}`))
	configFile.Close()
	return err
}
