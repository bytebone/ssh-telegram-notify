package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"justrainer/ssh-telegram-notify/tgsend"
	"net/http"
	"os"
	"strings"
	"time"

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
	log.Debugf("Reading Values")
	hostname := getHostname()
	loginDate := getLoginDate()
	user := getConnectedUser()
	loginIP := getLoginIP()

	log.Debug("Constructing Telegram Message")
	message := constructMessage(hostname, loginDate, user, loginIP)

	log.Debug("Sending Telegram Message")
	err := tgsend.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}
}

type ipLocation struct {
	Status        string `json:"status"`
	StatusMessage string `json:"message"`
	IP            string `json:"query"`
	Country       string `json:"country"`
	CountryCode   string `json:"countryCode"`
	City          string `json:"city"`
	ISP           string `json:"isp"`
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

func geoLocateIP(ip string) (ipLocation, error) {
	// send ip to http://ip-api.com/json/ and interpret returned json
	url := "http://ip-api.com/json/" + ip + "?fields=57875"
	response, err := http.Get(url)
	if err != nil {
		return ipLocation{}, err
	}
	if response.StatusCode != http.StatusOK {
		return ipLocation{}, fmt.Errorf("%s", response.Status)
	}
	defer response.Body.Close()

	var location ipLocation
	err = json.NewDecoder(response.Body).Decode(&location)
	if err != nil {
		return ipLocation{}, err
	}

	return location, nil
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
		ipLocation, err := geoLocateIP(ipAddress)
		if err != nil {
			// can be both http and json errors
			log.Warn(err)
			return
		}

		if ipLocation.Status != "success" {
			message += fmt.Sprintf("%v IP Address is %s\n", emoji.Warning, ipLocation.StatusMessage)
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
