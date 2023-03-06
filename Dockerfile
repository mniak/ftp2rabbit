FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
env GOOS=linux GOARCH=amd64

RUN go build -o main .

CMD ["/app/main"]

