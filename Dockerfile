FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o ./ftp2rabbit

FROM scratch
COPY --from=builder /app/ftp2rabbit /app/ftp2rabbit

ENTRYPOINT ["/app/ftp2rabbit"]
EXPOSE 10021

