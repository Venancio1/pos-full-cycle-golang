package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WeatherTempResponse struct {
	Current TempData `json:"current"`
}

type TempData struct {
	TempC  float64 `json:"temp_c"`
	TempF  float64 `json:"temp_f"`
	Kelvin float64
}

var ApiKey = "79362ad349514e759ea114534253011"
var WeatherApiBase = "https://api.weatherapi.com/v1/current.json"

func BuscaTemperatura(location string) (*TempData, error) {
	estado := url.QueryEscape(location)

	fullUrl := fmt.Sprintf("%s?key=%s&q=%s", WeatherApiBase, ApiKey, estado)
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	var weather WeatherTempResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	tempC := weather.Current.TempC
	tempF := weather.Current.TempF
	tempK := tempC + 273.15

	return &TempData{
		TempC:  tempC,
		TempF:  tempF,
		Kelvin: tempK,
	}, nil
}
