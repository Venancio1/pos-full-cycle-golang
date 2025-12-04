go-desafio-cloud-run
Explicação do Projeto
Weather Service API Este projeto é uma API em Go que fornece informações de temperatura baseadas em um código postal (CEP). A API utiliza uma integração com um serviço de clima e o serviço ViaCEP para converter CEPs em , retornando a temperatura atual em graus Celsius, Fahrenheit e Kelvin para uma determinada cidade.

Funcionalidades Consulta de cidade por CEP: Utiliza a API ViaCEP para converter um CEP em uma cidade. Consulta de temperatura por cidade: Utiliza a API WeatherAPI para obter a temperatura atual de uma cidade. Conversão de temperatura: Converte a temperatura de Celsius para Fahrenheit e Kelvin.

Endpoints
GET /weather?cep={CEP} Este endpoint recebe um CEP como parâmetro de query e retorna a temperatura na cidade correspondente.

Como executar a aplicação no Cloud Run:
https://cloudrun-fullcycle-go-848626908054.us-central1.run.app/infoclima?cep=06045270

Body de retorno

Response Body: { "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }