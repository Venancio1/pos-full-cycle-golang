# Weather Service API - Go Cloud Run

## 📋 Explicação do Projeto

Este projeto é uma API em Go que fornece informações de temperatura baseadas em um código postal (CEP). A API utiliza integração com um serviço de clima e o serviço ViaCEP para converter CEPs em localidades, retornando a temperatura atual em graus Celsius, Fahrenheit e Kelvin para uma determinada cidade.

---

## ✨ Funcionalidades

- **Consulta de cidade por CEP**: Utiliza a API ViaCEP para converter um CEP em uma cidade
- **Consulta de temperatura por cidade**: Utiliza a API WeatherAPI para obter a temperatura atual
- **Conversão de temperatura**: Converte a temperatura de Celsius para Fahrenheit e Kelvin

---

## 🔌 Endpoints

### GET `/weather?cep={CEP}`

Recebe um CEP como parâmetro de query e retorna a temperatura na cidade correspondente.

**Exemplo de requisição:**
```
https://cloudrun-fullcycle-go-848626908054.us-central1.run.app/infoclima?cep=06045270
```

**Exemplo de resposta:**
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

---

## 🚀 Como executar a aplicação no Cloud Run

A aplicação está disponível em:
```
https://cloudrun-fullcycle-go-848626908054.us-central1.run.app
```

---

## 📦 Estrutura do Projeto

```
.
├── client/                    # Clientes HTTP para APIs externas
│   ├── client_cep.go
│   ├── client_cep_test.go
│   ├── client_temperatura.go
│   └── client_temperatura_test.go
├── cmd/                       # Aplicação principal
│   ├── main.go
│   └── main_test.go
├── usecase/                   # Lógica de negócio
│   ├── busca_temperatura.go
│   └── busca_temperatura_test.go
├── Dockerfile                 # Configuração Docker
├── go.mod                     # Dependências do projeto
└── README.md                  # Este arquivo
```

---

## 🛠️ Desenvolvimento

### Pré-requisitos
- Go 1.21 ou superior
- Docker (para containerizar)

### Executar localmente
```bash
go run ./cmd/main.go
```

### Executar testes
```bash
go test ./...
```