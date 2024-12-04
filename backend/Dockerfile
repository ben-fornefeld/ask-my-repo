# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/server

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
