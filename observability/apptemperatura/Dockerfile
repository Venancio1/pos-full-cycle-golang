FROM golang:1.24 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apptemperatura ./cmd

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=build /app/apptemperatura .
ENTRYPOINT ["./apptemperatura"]