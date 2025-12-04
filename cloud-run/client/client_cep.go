package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var ViaCep = "https://viacep.com.br/ws/"

type AddressViaCep struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

func BuscaCep(cep string) *AddressViaCep {
	ctx := context.Background()
	url := fmt.Sprintf("%s%s/json/", ViaCep, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Erro na criação da requisição:", err.Error())
		return nil
	}
	println("Requisição criada com sucesso")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro na requisição:", err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var address AddressViaCep
	if err := json.Unmarshal(body, &address); err != nil {
		return nil
	}
	return &address
}
