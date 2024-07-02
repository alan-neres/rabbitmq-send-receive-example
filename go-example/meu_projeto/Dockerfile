# Use a imagem oficial do Go como base
FROM golang:1.17-alpine AS build

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o código fonte para o diretório de trabalho
COPY . .

# Baixe as dependências do Go
RUN go mod download

# Compile o aplicativo Go
RUN go build -o main .

# Crie uma imagem mínima do Alpine para executar o aplicativo compilado
FROM alpine:latest

# Defina o diretório de trabalho dentro do novo estágio da imagem
WORKDIR /root/

# Copie o executável do estágio de compilação para o diretório atual
COPY --from=build /app/main .

# Instale o cliente RabbitMQ (se necessário) - opcional
RUN apk add --no-cache rabbitmq-c

# Execute o aplicativo compilado
CMD ["./main"]

