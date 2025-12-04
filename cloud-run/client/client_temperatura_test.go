package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuscaTemperatura(t *testing.T) {
	t.Run("Temperatura válida retorna dados corretos", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(WeatherTempResponse{
				Current: TempData{TempC: 20.0, TempF: 68.0},
			})
		}))
		defer ts.Close()

		// Salva o valor original e redireciona para o servidor fake
		originalWeatherApiBase := WeatherApiBase
		WeatherApiBase = ts.URL
		defer func() { WeatherApiBase = originalWeatherApiBase }()

		temp, err := BuscaTemperatura("SP")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if temp.TempC != 20.0 {
			t.Errorf("esperava 20.0°C, veio %.2f", temp.TempC)
		}
		if temp.TempF != 68.0 {
			t.Errorf("esperava 68.0°F, veio %.2f", temp.TempF)
		}
		if temp.Kelvin != 293.15 {
			t.Errorf("esperava 293.15K, veio %.2f", temp.Kelvin)
		}
	})

	t.Run("JSON inválido retorna erro", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer ts.Close()

		originalWeatherApiBase := WeatherApiBase
		WeatherApiBase = ts.URL
		defer func() { WeatherApiBase = originalWeatherApiBase }()

		_, err := BuscaTemperatura("SP")
		if err == nil {
			t.Fatalf("esperava erro ao decodificar JSON, veio nil")
		}
	})

	t.Run("Temperatura em diferentes unidades", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(WeatherTempResponse{
				Current: TempData{TempC: 25.0, TempF: 77.0},
			})
		}))
		defer ts.Close()

		originalWeatherApiBase := WeatherApiBase
		WeatherApiBase = ts.URL
		defer func() { WeatherApiBase = originalWeatherApiBase }()

		temp, err := BuscaTemperatura("RJ")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if temp.TempC != 25.0 {
			t.Errorf("esperava 25.0°C, veio %.2f", temp.TempC)
		}
		if temp.TempF != 77.0 {
			t.Errorf("esperava 77.0°F, veio %.2f", temp.TempF)
		}
		if temp.Kelvin != 298.15 {
			t.Errorf("esperava 298.15K, veio %.2f", temp.Kelvin)
		}
	})
}
