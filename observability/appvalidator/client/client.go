package client

import (
	"appvalidator/otel"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/propagation"
)

type Cep struct {
	CEP string `json:"cep"`
}

type ResponseClima struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

const url = "http://apptemperatura:8080/infoclima"

func (c *Cep) Validate() error {
	if len(c.CEP) != 8 {
		return errors.New("invalid zipcode: length must be 8")
	}
	return nil
}

func (c *Cep) GetWeatherWithContext(ctx context.Context) (*ResponseClima, error) {
	tr := otel.Tracer("GetWeatherWithContext")
	ctx, span := tr.Start(ctx, "CallAppTemperaturaAPI")
	defer span.End()

	body, err := json.Marshal(c)
	if err != nil {
		log.Println("Erro ao serializar CEP:", err.Error())
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		log.Println("Erro na criação da requisição:", err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Inject trace context into headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro na requisição:", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ResponseClima
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
