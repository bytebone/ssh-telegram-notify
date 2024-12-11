package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IpLocation struct {
	Status      int     `json:",omitempty"`
	IP          string  `json:"ip"`
	Country     string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	ISP         string  `json:"org"`
	Latitude    float32 `json:"latitude"`
	Longitude   float32 `json:"longitude"`
}

func Locate(ip string) (IpLocation, error) {
	url := "https://ipapi.co/" + ip + "/json/"
	response, err := http.Get(url)
	if err != nil {
		return IpLocation{Status: response.StatusCode}, err
	}
	if response.StatusCode != 200 {
		return IpLocation{Status: response.StatusCode}, fmt.Errorf("%s", response.Status)
	}
	defer response.Body.Close()

	var decodedLocation IpLocation
	err = json.NewDecoder(response.Body).Decode(&decodedLocation)
	if err != nil {
		return IpLocation{}, err
	}

	decodedLocation.Status = response.StatusCode
	return decodedLocation, nil
}
