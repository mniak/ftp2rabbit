FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux GOARCH=amd64

RUN go build -o main .

CMD ["/app/main"]