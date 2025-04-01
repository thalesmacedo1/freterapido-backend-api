# API de Cotação de Frete - Frete Rápido

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-Framework-brightgreen)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue)
![GORM](https://img.shields.io/badge/GORM-ORM-lightblue)
![DDD](https://img.shields.io/badge/DDD-Architecture-orange)
![Clean Architecture](https://img.shields.io/badge/Clean-Architecture-red)

## Sumário

- [Descrição](#descrição)
- [Arquitetura](#arquitetura)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação e Configuração](#instalação-e-configuração)
- [Rotas da API](#rotas-da-api)
- [Testes](#testes)
- [Docker e Docker Compose](#docker-e-docker-compose)

## Descrição

A API de Cotação de Frete é um sistema para consulta de valores de frete através de integrações com transportadoras. Desenvolvida com Go, segue os princípios de Domain-Driven Design (DDD) e Clean Architecture, proporcionando uma solução modular, testável e de fácil manutenção.

### Funcionalidades Principais

- Consulta de cotações de frete com diferentes transportadoras
- Métricas e análises estatísticas sobre as cotações realizadas
- Armazenamento de histórico de cotações para consulta posterior

## Arquitetura

O projeto segue uma arquitetura em camadas baseada nos princípios de DDD e Clean Architecture:

```
api/
├── domain/                  # Camada de Domínio
│   └── entities/            # Entidades e regras de negócio
│       
├── application/             # Camada de Aplicação
│   └── usecases/            # Casos de uso da aplicação
│       
├── infrastructure/          # Camada de Infraestrutura
│   └── database/            # Implementações de persistência
│       
├── interfaces/              # Camada de Interface
│   ├── api/                 # Controllers HTTP
│   └── routers/             # Configuração de rotas
│       
└── cmd/                     # Pontos de entrada do sistema
    └── api/                 # Servidor HTTP
```

### Benefícios da Arquitetura

- **Desacoplamento**: Cada camada possui responsabilidades bem definidas
- **Testabilidade**: Facilidade para escrever testes unitários e de integração
- **Manutenibilidade**: Código organizado e de fácil compreensão
- **Escalabilidade**: Facilidade para adicionar novos recursos

## Tecnologias Utilizadas

- **Linguagem**: Go 1.23+
- **Framework HTTP**: Gin
- **ORM**: GORM
- **Banco de Dados**: PostgreSQL
- **Testes**: Testify, Go-SQLMock
- **Containerização**: Docker e Docker Compose

## Instalação e Configuração

### Pré-requisitos

- Go 1.23+
- PostgreSQL
- Docker e Docker Compose (opcional)

### Instalação local

1. Clone o repositório:
   ```bash
   git clone https://github.com/thalesmacedo1/freterapido-backend-api.git
   cd freterapido-backend-api
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Configure o banco de dados:
   - Certifique-se de que o PostgreSQL está em execução
   - Crie um banco de dados para a aplicação
   - Configure as variáveis de ambiente ou ajuste a string de conexão em `api/cmd/api/main.go`

4. Execute a aplicação:
   ```bash
   go run api/cmd/api/main.go
   ```

### Instalação com Docker

1. Clone o repositório:
   ```bash
   git clone https://github.com/thalesmacedo1/freterapido-backend-api.git
   cd freterapido-backend-api
   ```

2. Inicie os containers:
   ```bash
   docker-compose up -d
   ```

## Rotas da API

A API disponibiliza os seguintes endpoints:

### 1. Cotação de Frete

**Endpoint**: `POST /quote`

**Descrição**: Recebe dados do destinatário e volumes para realizar cotação de frete com diferentes transportadoras.

**Requisição**:
```json
{
  "recipient": {
    "address": {
      "zipcode": "01311000"
    }
  },
  "volumes": [
    {
      "category": 7,
      "amount": 1,
      "unitary_weight": 5,
      "price": 349,
      "sku": "abc-teste-123",
      "height": 0.2,
      "width": 0.2,
      "length": 0.2
    },
    {
      "category": 7,
      "amount": 2,
      "unitary_weight": 4,
      "price": 556,
      "sku": "abc-teste-527",
      "height": 0.4,
      "width": 0.6,
      "length": 0.15
    }
  ]
}
```

**Resposta**:
```json
{
  "carrier": [
    {
      "name": "EXPRESSO FR",
      "service": "Rodoviário",
      "deadline": "3",
      "price": 17
    },
    {
      "name": "Correios",
      "service": "SEDEX",
      "deadline": "1",
      "price": 20.99
    }
  ]
}
```

### 2. Métricas de Cotações

**Endpoint**: `GET /metrics?last_quotes={quantidade}`

**Descrição**: Retorna métricas sobre as cotações realizadas. O parâmetro `last_quotes` é opcional e limita a análise às N cotações mais recentes.

**Parâmetros de consulta**:
- `last_quotes` (opcional): Número inteiro que indica quantas cotações recentes devem ser consideradas na análise

**Resposta**:
```json
{
  "carrier_metrics": [
    {
      "carrier_name": "EXPRESSO FR",
      "total_quotes": 10,
      "total_shipping_price": 150.50,
      "average_shipping_price": 15.05
    },
    {
      "carrier_name": "Correios",
      "total_quotes": 5,
      "total_shipping_price": 120.25,
      "average_shipping_price": 24.05
    }
  ],
  "cheapest_and_most_expensive": {
    "cheapest_shipping": 12.50,
    "most_expensive_shipping": 30.75
  }
}
```

## Testes

O projeto inclui testes unitários e de integração para validar o funcionamento correto da aplicação.

### Executando os testes

#### Testes unitários
```bash
make test-unit
```

#### Testes de integração
```bash
make test-integration
```

#### Todos os testes
```bash
make test
```

## Docker e Docker Compose

O projeto inclui arquivos para execução em ambiente Docker:

- `Dockerfile`: Configuração para construir a imagem da aplicação
- `docker-compose.yml`: Configuração para orquestrar os containers da aplicação e banco de dados

### Comandos úteis

- Iniciar aplicação e banco de dados:
  ```bash
  make docker-up
  ```

- Parar containers:
  ```bash
  make docker-down
  ```

- Construir aplicação:
  ```bash
  make build
  ```

- Executar aplicação localmente:
  ```bash
  make run
  ```

## Observações sobre a API do Frete Rápido

Para consumir a API do Frete Rápido, os seguintes dados são obrigatórios:

- CNPJ Remetente: 25.438.296/0001-58 (mesmo valor para "shipper.registered_number" e "dispatchers.registered_number")
- Token de autenticação: 1d52a9b6b78cf07b08586152459a5c90
- Código Plataforma: 5AKVkHqCn
- CEP: 29161-376 (dispatchers[*].zipcode)
- O campo "unitary_price" deve ser informado