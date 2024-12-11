package telegram

import (
	"bytebone/ssh-telegram-notify/internal/config"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

func SendMessage(config *config.MessageBackendTelegram, params *SendMessageParams) error {
	switch {
	case config.Token == "":
		return errors.New("no telegram bot token found")
	case config.ChatID == "":
		return errors.New("no telegram chat id found")
	case config.Token == "" && config.ChatID == "":
		return errors.New("no telegram bot token or group id found")
	}

	data := url.Values{
		"chat_id":      {config.ChatID},
		"text":         {params.Message},
		"parse_mode":   {"Markdown"},
		"reply_markup": {fmt.Sprintf(`{"inline_keyboard":[[{"text": "Show on Shodan","url":"https://www.shodan.io/host/%s"}]]}`, params.IP)},
	}
	if config.MessageThreadID != "" {
		data.Add("message_thread_id", config.MessageThreadID)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.Token)
	res, err := http.PostForm(url, data)
	log.Debug(url)
	log.Debug(data)
	if res.StatusCode != 200 {
		return fmt.Errorf("%s", res.Status)
	}
	return err
}

type SendMessageParams struct {
	Message string
	IP      string
	// Location IPLocation
}

type IPLocation struct {
	Latitude  float32
	Longitude float32
}
