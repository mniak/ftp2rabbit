FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOOS=linux GOARCH=amd64
ENV PATH=/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN go build -o main .

# EXPOSE 20-21/tcp 10021/tcp 10000-10020/tcp

CMD ["/app/main"]