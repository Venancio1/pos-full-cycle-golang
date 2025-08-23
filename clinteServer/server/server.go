package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func main() {
	log.Println("Starting server...")

	var err error
	Db, err = NewDB()
	if err != nil {
		log.Println("Erro na conexão com o banco de dados:", err.Error())
		return
	}

	Server()
}

type CotacaoDb struct {
	ID           uint
	CotacaoDolar string
}

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func Server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoDolarHandler)
	http.ListenAndServe(":8080", mux)
}

func cotacaoDolarHandler(w http.ResponseWriter, r *http.Request) {
	response, err := BuscaCotacao()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func BuscaCotacao() (*string, error) {
	ctx, cancel := timeoutCtx(200 * time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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
	var c Cotacao
	if err := json.Unmarshal(body, &c); err != nil {
		return nil, err
	}

	var cdb CotacaoDb
	cdb.CotacaoDolar = c.USDBRL.Bid

	ctx_db, cancel := timeoutCtx(10 * time.Millisecond)
	defer cancel()
	if err := SaveCotacao(ctx_db, Db, &cdb); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Tempo de requisição excedido")
			return nil, err
		}
		return nil, err
	}

	return &c.USDBRL.Bid, nil
}

func NewDB() (*sql.DB, error) {
	log.Println("Iniciando banco de dados...")
	db, err := sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		return nil, err
	}
	createTable := `
    CREATE TABLE IF NOT EXISTS cotacoes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        cotacao TEXT NOT NULL
    );`

	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}
	return db, nil
}

func SaveCotacao(ctx context.Context, db *sql.DB, cotacao *CotacaoDb) error {
	query := `INSERT INTO cotacoes (cotacao) VALUES (?)`
	_, err := db.Exec(query, cotacao.CotacaoDolar)
	return err
}

func timeoutCtx(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
