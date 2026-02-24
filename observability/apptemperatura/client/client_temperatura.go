package client

import (
	"apptemperatura/otel"
	"context"
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
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

var ApiKey = "79362ad349514e759ea114534253011"
var WeatherApiBase = "https://api.weatherapi.com/v1/current.json"

func BuscaTemperatura(ctx context.Context, location string) (*TempData, error) {
	tr := otel.Tracer("BuscaTemperatura")
	ctx, span := tr.Start(ctx, "BuscaTemperaturaAPI")
	defer span.End()

	estado := url.QueryEscape(location)

	fullUrl := fmt.Sprintf("%s?key=%s&q=%s", WeatherApiBase, ApiKey, estado)
	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
	if err != nil {
		log.Println("Erro na criação da requisição de temperatura:", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro na requisição de temperatura:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erro ao ler resposta de temperatura:", err.Error())
		return nil, err
	}

	fmt.Println(string(body))
	var weather WeatherTempResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON: %w", err)
	}

	tempC := weather.Current.Temp_C
	tempF := weather.Current.Temp_F
	tempK := tempC + 273.15

	return &TempData{
		Temp_C: tempC,
		Temp_F: tempF,
		Temp_K: tempK,
	}, nil
}
