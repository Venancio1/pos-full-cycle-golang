# Observability - Arquitetura de Microserviços com Tracing Distribuído

Projeto demonstrando uma arquitetura de microserviços em Go com observabilidade completa através de rastreamento distribuído (OpenTelemetry + Zipkin).

# Subir todos os serviços (com rebuild das imagens)
docker compose up -d --build


## 🧪 Exemplos de Testes com curl

```bash
# Requisição simples
curl -X POST http://localhost:8082/cep \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}'

# Resposta esperada:
# {"city":"São Paulo","temp_c":21,"temp_f":69.9,"temp_k":294.15}


Após executar os testes com curl, você pode visualizar os traces:

1. Abra seu navegador em: `http://localhost:9411`
2. Clique em "Run Query" para ver os serviços disponíveis
3. Selecione um serviço (apptemperatura ou appvalidator) da dropdown
4. Clique em "Run Query" novamente
5. Clique em um trace para ver os detalhes das chamadas, incluindo:
   - Tempo de execução de cada span
   - Tempo de resposta das APIs externas
   - Propagação de contexto entre serviços


## 📝 Variáveis de Ambiente

### apptemperatura
```
PORT=8080  # Porta do serviço (padrão: 8080)
```

### appvalidator
```
PORT=8080  # Porta do serviço (padrão: 8080)
```

### OpenTelemetry
```
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317  # Configurado nos apps
```
