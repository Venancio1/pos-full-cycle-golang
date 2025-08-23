package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	println("Starting client...")
	Client()
}

type Cotacao struct {
	Dolar string `json:"bid"`
}

func Client() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println("Erro na criação da requisição:", err.Error())
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Tempo de requisição excedido")
			return
		}
		log.Println("Erro na requisição:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println("Erro ao ler o body da resposta", err.Error())
		return
	}

	var c Cotacao
	err = json.Unmarshal(body, &c.Dolar)
	if err != nil {
		println("Erro ao desserializar a resposta:", err.Error())
		return
	}

	err = GeraArquivo(c)
	if err != nil {
		println("Erro ao gerar arquivo:", err.Error())
		return
	}
}

func GeraArquivo(cotacao Cotacao) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		println("erro ao criar arquivo:")
		return err
	}
	defer file.Close()
	_, err = file.WriteString(("Dólar: " + cotacao.Dolar))
	if err != nil {
		println("erro ao escrever no arquivo:", err.Error())
		return err
	}
	return nil
}
