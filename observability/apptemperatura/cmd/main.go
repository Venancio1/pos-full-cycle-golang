package main

import (
	"apptemperatura/client"
	"apptemperatura/otel"
	"apptemperatura/usecase"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"go.opentelemetry.io/otel/propagation"
)

func main() {
	Server()
}

func Server() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/infoclima", InfoClimaHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback para rodar local
	}
	http.ListenAndServe(":"+port, mux)
}

func InfoClimaHandler(w http.ResponseWriter, r *http.Request) {
	// Extract propagated context and start a span
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tr := otel.Tracer("apptemperatura")
	ctx, span := tr.Start(ctx, "InfoClima")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var cepDt client.CepData
	if err := json.Unmarshal(body, &cepDt); err != nil {
		http.Error(w, "Erro ao deserializar cep", http.StatusBadRequest)
		return
	}

	println("CEP válido:", cepDt.CEP)

	response, err := usecase.BuscaClimatizacao(ctx, cepDt.CEP)
	if err != nil {
		if strings.HasPrefix(err.Error(), "notfound") {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}

		http.Error(w, "Erro ao buscar informações de clima", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Erro ao serializar resposta", http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
