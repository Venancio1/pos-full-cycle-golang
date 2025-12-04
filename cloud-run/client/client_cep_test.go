package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuscaCep(t *testing.T) {
	t.Run("CEP válido retorna estado correto", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(AddressViaCep{Localidade: "Carapicuiba", Erro: false})
		}))
		defer ts.Close()

		// redireciona ViaCep para o servidor fake
		originalViaCep := ViaCep
		ViaCep = ts.URL + "/"
		defer func() { ViaCep = originalViaCep }()

		address := BuscaCep("")
		if address == nil {
			t.Fatalf("esperava endereço, mas veio nil")
		}
		if address.Localidade != "Carapicuiba" {
			t.Errorf("esperava , mas veio %s", address.Localidade)
		}
		if address.Erro {
			t.Errorf("esperava Erro=false, mas veio true")
		}
	})

	t.Run("CEP inválido retorna erro", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(AddressViaCep{Localidade: "", Erro: true})
		}))
		defer ts.Close()

		originalViaCep := ViaCep
		ViaCep = ts.URL + "/"
		defer func() { ViaCep = originalViaCep }()

		address := BuscaCep("06341580")
		if address == nil {
			t.Fatalf("esperava endereço, mas veio nil")
		}
		if !address.Erro {
			t.Errorf("esperava Erro=true, mas veio false")
		}
	})

	t.Run("Falha na requisição retorna nil", func(t *testing.T) {
		// redireciona para URL inválida
		originalViaCep := ViaCep
		ViaCep = "http://invalid-url-that-does-not-exist/"
		defer func() { ViaCep = originalViaCep }()

		address := BuscaCep("06341580")
		if address != nil {
			t.Errorf("esperava nil em caso de erro, mas veio %v", address)
		}
	})
}
