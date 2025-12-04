package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInfoClimaHandler_MissingCep(t *testing.T) {
	req := httptest.NewRequest("GET", "/infoclima", nil)
	w := httptest.NewRecorder()

	InfoClimaHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("esperava 400, mas veio %d", w.Code)
	}
}

func TestValidationCep(t *testing.T) {
	if err := validationCep("12345678"); err != nil {
		t.Errorf("CEP válido deveria passar, mas deu erro: %v", err)
	}

	if err := validationCep("123"); err == nil {
		t.Errorf("CEP inválido deveria falhar")
	}
}
