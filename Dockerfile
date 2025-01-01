# Stage 1: Build
FROM golang:1.23.4-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add build-base

RUN go install github.com/a-h/templ/cmd/templ@v0.2.793

RUN go mod download

COPY . .

RUN CGO_ENABLED='1' templ generate .
RUN CGO_ENABLED='1' go build cmd/main.go

# Stage 2: Minimal Runtime
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main trends

CMD ["./trends"]