package usecase

import (
	"cloudrun/client"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuscaClimatizacao(t *testing.T) {
	t.Run("CEP válido retorna clima correto", func(t *testing.T) {
		// mock do servidor de CEP
		cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(client.AddressViaCep{Localidade: "Osasco", Erro: false})
		}))
		defer cepServer.Close()

		// mock do servidor de temperatura
		tempServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(client.WeatherTempResponse{
				Current: client.TempData{TempC: 25.0, TempF: 77.0},
			})
		}))
		defer tempServer.Close()

		// Salva os valores originais
		originalViaCep := client.ViaCep
		originalWeatherApiBase := client.WeatherApiBase

		// Redireciona para os servidores fake
		client.ViaCep = cepServer.URL + "/"
		client.WeatherApiBase = tempServer.URL

		defer func() {
			client.ViaCep = originalViaCep
			client.WeatherApiBase = originalWeatherApiBase
		}()

		resp, err := BuscaClimatizacao("01001000")
		if err != nil {
			t.Fatalf("esperava sucesso, mas deu erro: %v", err)
		}
		if resp.TempC != 25.0 {
			t.Errorf("esperava 25.0°C, mas veio %.2f", resp.TempC)
		}
		if resp.TempF != 77.0 {
			t.Errorf("esperava 77.0°F, mas veio %.2f", resp.TempF)
		}
		if resp.Kelvin != 298.15 {
			t.Errorf("esperava 298.15K, mas veio %.2f", resp.Kelvin)
		}
	})

	t.Run("CEP inválido retorna erro", func(t *testing.T) {
		// mock do servidor de CEP com erro
		cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(client.AddressViaCep{Localidade: "", Erro: true})
		}))
		defer cepServer.Close()

		originalViaCep := client.ViaCep
		client.ViaCep = cepServer.URL + "/"

		defer func() {
			client.ViaCep = originalViaCep
		}()

		resp, err := BuscaClimatizacao("00000000")
		if err == nil {
			t.Fatalf("esperava erro para CEP inválido, mas deu sucesso: %v", resp)
		}
	})
}
