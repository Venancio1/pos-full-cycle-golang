package main

import (
	"appvalidator/client"
	"appvalidator/otel"
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
	mux.HandleFunc("/cep", CepHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback para rodar local
	}
	http.ListenAndServe(":"+port, mux)
}

func CepHandler(w http.ResponseWriter, r *http.Request) {
	// Extract propagated context and start a span
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	tr := otel.Tracer("appvalidator")
	ctx, span := tr.Start(ctx, "CepValidation")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}
	if string(body) == "" {
		http.Error(w, "cep é obrigatório", http.StatusBadRequest)
		return
	}
	var c client.Cep
	err = json.Unmarshal(body, &c)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}

	err = c.Validate()
	if err != nil {
		http.Error(w, "invalid zipcode: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	response, err := c.GetWeatherWithContext(ctx)
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
