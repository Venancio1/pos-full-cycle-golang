package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type AddressBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type AddressViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const (
	UrlBrasilApi = "https://brasilapi.com.br/api/cep/v1/01001000"
	UrlViaCep    = "https://viacep.com.br/ws/01001000/json/"
)

func main() {

	var brApiCh = make(chan interface{})
	var viaCepCh = make(chan interface{})

	go func() {
		result, err := BuscaCep(UrlBrasilApi)
		if err != nil {
			brApiCh <- err
			log.Println("Erro na busca do cep no endereço Brasil API:", err.Error())
			return
		}
		brApiCh <- result
	}()
	go func() {
		result, err := BuscaCep(UrlViaCep)
		if err != nil {
			viaCepCh <- err
			log.Println("Erro na busca do cep no endereço Via Cep:", err.Error())
			return
		}
		viaCepCh <- result
	}()

	select {
	case res1 := <-brApiCh:
		fmt.Println(UrlBrasilApi, res1)
	case res2 := <-viaCepCh:
		fmt.Println(UrlViaCep, res2)
	}
}

func BuscaCep(url string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("Erro na criação da requisição:", err.Error())
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Tempo de requisição excedido")
			return nil, err
		}
		log.Println("Erro na requisição:", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch {
	case strings.Contains(url, "brasilapi"):
		var address AddressBrasilApi
		if err := json.Unmarshal(body, &address); err != nil {
			return nil, err
		}
		return address, nil
	case strings.Contains(url, "viacep"):
		var address AddressViaCep
		if err := json.Unmarshal(body, &address); err != nil {
			return nil, err
		}
		return address, nil
	}

	return nil, errors.New("URL desconhecida")
}
