# Use uma imagem Go base
FROM golang:1.21-alpine as builder

# Setando o diretório de trabalho
WORKDIR /app

# Copiar os arquivos do Go para o container
COPY . .

# Baixar dependências e compilar o projeto
RUN go mod tidy
RUN go build -o client ./client.go

# Rodar a aplicação
CMD ["./client"]
