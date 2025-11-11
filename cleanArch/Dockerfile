FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o ordersystem ./cmd/ordersystem

EXPOSE 3001

CMD ["./ordersystem"]