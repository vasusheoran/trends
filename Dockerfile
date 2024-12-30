# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/trends -v cmd/main.go

# Stage 2: Minimal Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/trends . 

CMD ["./trends"]