package ip

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IpLocation struct {
	Status      int    `json:",omitempty"`
	IP          string `json:"ip"`
	Country     string `json:"country_name"`
	CountryCode string `json:"country_code"`
	City        string `json:"city"`
	ISP         string `json:"org"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

func GeoLocateIP(ip string) (IpLocation, error) {
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

func DrawIPMap(lat int, lon int) (locationMapImageURL string, err error) {
	url := fmt.Sprintf("https://maps.geoapify.com/v1/staticmap?style=osm-carto&width=600&height=400&center=lonlat:%d,%d&zoom=5&marker=lonlat:%d,%d;color:%%23ff0000;size:small&apiKey=2d94c1faf1304ec090791ed6584a1991", lon, lat, lon, lat)
	return url, nil
}
