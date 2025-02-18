# Usando uma versão mais recente do Go
FROM golang:1.23 AS builder

RUN apt-get update && apt-get install -y ca-certificates

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos do projeto para dentro do container
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte restante
COPY . .

# Compila o aplicativo
RUN go build -o main .

# Usa uma imagem mínima para rodar o binário gerado
FROM debian:bookworm-slim

# Define o diretório de trabalho no runtime
WORKDIR /root/

# Copia o binário do build para a imagem final
COPY --from=builder /app/main .

# Expõe a porta da aplicação
EXPOSE 8080

# Comando padrão para rodar a aplicação
CMD ["./main"]
