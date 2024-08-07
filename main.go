package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"bytebone/ssh-telegram-notify/internal/config"
	"bytebone/ssh-telegram-notify/internal/ip"
	"bytebone/ssh-telegram-notify/internal/senders/telegram"
	"bytebone/ssh-telegram-notify/internal/utils"

	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
)

func init() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

func main() {
	var configPath = utils.GetHomeDir() + "/.config/ssh-notify/config.json"
	config.DecodeConfig(configPath)

	log.Debugf("Reading Values")
	hostname := getHostname()
	loginDate := getLoginDate()
	user := getConnectedUser()
	loginIP := getLoginIP()

	log.Debug("Constructing Telegram Message")
	message := constructMessage(hostname, loginDate, user, loginIP)

	log.Debug("Sending Telegram Message")
	err := telegram.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}
}

func getHostname() (hostname string) {
	hostname, err := os.Hostname()
	if err != nil {
		log.Warn(err)
	}
	return
}

func getLoginDate() (loginDate string) {
	now := time.Now()
	loginDate = now.Format("02.01.2006 15:04:05")
	return
}

func getConnectedUser() (user string) {
	user = os.Getenv("USER")
	return
}

func getLoginIP() (loginIP string) {
	loginIP = strings.Split(os.Getenv("SSH_CLIENT"), " ")[0]
	return
}

func constructMessage(hostname string, loginDate string, user string, ipAddress string) (message string) {
	message = fmt.Sprintf("*Session started on %s*\n\n", hostname)

	if loginDate != "" {
		message += fmt.Sprintf("%v %s\n", emoji.Calendar, loginDate)
	}

	if user != "" {
		message += fmt.Sprintf("%v %s\n", emoji.BustInSilhouette, user)
	}

	if ipAddress != "" {
		ipLocation, err := ip.GeoLocateIP(ipAddress)
		if err != nil {
			// can be both http and json errors
			log.Warn(err)
			return
		}

		if ipLocation.Status != 200 {
			message += fmt.Sprintf("%v Couldn't resolve IP\n", emoji.Warning)
			return
		}

		countryEmoji, err := emoji.CountryFlag(ipLocation.CountryCode)
		if err != nil {
			countryEmoji = emoji.GlobeShowingEuropeAfrica
		}
		message += fmt.Sprintf("%v %s, %s\n", countryEmoji, ipLocation.City, ipLocation.Country)
		message += fmt.Sprintf("%v %s\n", emoji.ElectricPlug, ipLocation.ISP)
		message += fmt.Sprintf("%v `%s`\n", emoji.DesktopComputer, ipLocation.IP)
	}

	return
}
