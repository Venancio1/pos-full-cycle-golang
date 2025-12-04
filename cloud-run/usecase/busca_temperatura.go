package usecase

import (
	"cloudrun/client"
	"fmt"
)

type ResponseClima struct {
	TempC  float64
	TempF  float64
	Kelvin float64
}

func BuscaClimatizacao(cep string) (*ResponseClima, error) {

	address := client.BuscaCep(cep)
	if address == nil || address.Erro {
		return nil, fmt.Errorf("notfound: CEP não encontrado")
	}

	println("Endereço encontrado:", address.Localidade)

	temperature, err := client.BuscaTemperatura(address.Localidade)
	if err != nil {
		return nil, err
	}

	response := &ResponseClima{
		TempC:  temperature.TempC,
		TempF:  temperature.TempF,
		Kelvin: temperature.Kelvin,
	}

	return response, nil

}
