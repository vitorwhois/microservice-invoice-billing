FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário gerado a partir do estágio builder
COPY --from=builder /app/main .

# Se as migrações forem necessárias, copia a pasta de migrações
COPY migrations ./migrations

EXPOSE 8080

CMD ["./main"]
