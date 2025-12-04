package main

import (
	"cloudrun/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

func main() {
	Server()
}

func Server() {
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
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "Parâmetro 'cep' é obrigatório", http.StatusBadRequest)
		return
	}
	println("CEP recebido:", cep)

	err := validationCep(cep)
	if err != nil {
		http.Error(w, "invalid zipcode:", http.StatusUnprocessableEntity)
		return
	}

	println("CEP válido:", cep)
	response, err := usecase.BuscaClimatizacao(cep)
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

func validationCep(cep string) error {
	if len(cep) != 8 {
		return errors.New("cep inválido")
	}
	return nil
}
