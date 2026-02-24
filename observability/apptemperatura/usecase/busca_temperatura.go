package usecase

import (
	"apptemperatura/client"
	"apptemperatura/otel"
	"context"
	"fmt"
)

type ResponseClima struct {
	City   string  `json:"city"`
	Temp_C float64 `json:"temp_c"`
	Temp_F float64 `json:"temp_f"`
	Temp_K float64 `json:"temp_k"`
}

func BuscaClimatizacao(ctx context.Context, cep string) (*ResponseClima, error) {
	tr := otel.Tracer("BuscaClimatizacao")
	ctx, span := tr.Start(ctx, "SearchTemperature")
	defer span.End()

	address := client.BuscaCep(ctx, cep)
	if address == nil || address.Erro {
		return nil, fmt.Errorf("notfound: CEP não encontrado")
	}

	println("Endereço encontrado:", address.Localidade)

	temperature, err := client.BuscaTemperatura(ctx, address.Localidade)
	if err != nil {
		return nil, err
	}

	response := &ResponseClima{
		City:   address.Localidade,
		Temp_C: temperature.Temp_C,
		Temp_F: temperature.Temp_F,
		Temp_K: temperature.Temp_K,
	}

	return response, nil

}
