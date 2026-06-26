# Estágio de compilação (Build stage)
FROM golang:1.26-alpine AS builder

# Instalar dependências necessárias para compilação
RUN apk add --no-cache git

WORKDIR /app

# Copiar arquivos de dependências e baixar
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo o código-fonte
COPY . .

# Compilar o binário estaticamente
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bitbank .

# Estágio final (Final stage)
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar o binário compilado do estágio anterior
COPY --from=builder /app/bitbank .

# Porta padrão exposta
EXPOSE 8080

# Definir ponto de entrada padrão para iniciar a API REST
ENTRYPOINT ["/app/bitbank", "serve", "--port", "8080"]
