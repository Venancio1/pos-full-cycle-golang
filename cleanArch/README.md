📦 Projeto de Orders com Migrations no Docker
Este projeto em Go exemplifica a criação e gerenciamento de orders com uma arquitetura moderna e escalável. Ele integra:

Banco de dados MySQL

RabbitMQ para mensageria

Endpoint REST

Servidor gRPC

Servidor GraphQL

Migrações automatizadas via Docker

🚀 Desafio de Listagem de Orders
Passos para executar o projeto

1. Clone o repositório
2. Navegue até o diretório do projeto: cd nome-do-projeto
3. Monte e suba os serviços: docker-compose up --build
4. Caso o build já tenha sido feito anteriormente: docker-compose up
5. Se a migração não for executada automaticamente: docker-compose run migrate

📡 Serviços e Portas
REST API: http://localhost:8000/order
gRPC Server: Porta 50051
GraphQL API: http://localhost:8080/