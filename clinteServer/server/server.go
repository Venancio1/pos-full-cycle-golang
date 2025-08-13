package main

import (
	"encoding/json"
	"io"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func main() {
	println("Starting server...")
	Server()
}

type Cotacoes struct {
	CotacaoDolar string
}

type Bid struct {
	Bid string `json:"bid"`
}

type Cotacao struct {
	USDBRL Bid `json:"USDBRL"`
}

func (c *Cotacao) ToModel() *Cotacoes {
	return &Cotacoes{
		CotacaoDolar: c.USDBRL.Bid,
	}
}

func Server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoDolarHandler)
	http.ListenAndServe(":8080", mux)
}

func cotacaoDolarHandler(w http.ResponseWriter, r *http.Request) {
	response, err := buscaCotacao()
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

func buscaCotacao() (*Bid, error) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
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

	dataBase, err := NewDB()
	if err != nil {
		println("Erro na conexão com o banco de dados:", err.Error())
		return nil, err
	}
	if err := saveCotacao(dataBase, c.ToModel()); err != nil {
		return nil, err
	}
	return &c.USDBRL, nil
}

func NewDB() (*gorm.DB, error) {
	// Config padrão "modernc.org/sqlite" para não depender de compilador C.
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "sqlite",
		DSN:        "cotacoes.db",
	}), &gorm.Config{}) // db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Cotacoes{}); err != nil {
		return nil, err
	}
	return db, nil
}

func saveCotacao(db *gorm.DB, cotacao *Cotacoes) error {
	if err := db.Create(cotacao).Error; err != nil {
		return err
	}
	return nil
}
