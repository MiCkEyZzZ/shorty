# Сборка бинарника
FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o shorty ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/shorty .
COPY .env .env

EXPOSE 3000

CMD [ "./shorty" ]
