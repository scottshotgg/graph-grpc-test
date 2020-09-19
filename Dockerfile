FROM golang:1.15.2-alpine3.12 as builder

WORKDIR /app

COPY . .

RUN go build -o /bin/app ./cmd/exchange/main.go