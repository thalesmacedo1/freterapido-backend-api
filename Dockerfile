FROM golang:1.23-alpine

# Definindo o diretório de trabalho dentro do container
WORKDIR /app

# Copiando os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixando as dependências
RUN go mod download

# Copiando o código-fonte
COPY . .

# Definindo o diretório de trabalho gerar o build
WORKDIR /app/cmd/server

# Construindo o executável
RUN go build -o api

# Expondo a porta
EXPOSE 8080

# Comando para rodar a API
CMD ["./api"]