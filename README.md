# API de Cotação de Frete - Frete Rápido

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Gin](https://img.shields.io/badge/Gin-Framework-brightgreen)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue)
![GORM](https://img.shields.io/badge/GORM-ORM-lightblue)
![DDD](https://img.shields.io/badge/DDD-Architecture-orange)
![Clean Architecture](https://img.shields.io/badge/Clean-Architecture-red)
![Swagger](https://img.shields.io/badge/Swagger-Documentation-green)

## Sumário

- [Descrição](#descrição)
- [Arquitetura](#arquitetura)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação e Configuração](#instalação-e-configuração)
- [Rotas da API](#rotas-da-api)
- [Documentação Swagger](#documentação-swagger)

## Descrição

A API de Cotação de Frete é um sistema para consulta de valores de frete através de integrações com transportadoras. Desenvolvida com Go, segue os princípios de Domain-Driven Design (DDD) e Clean Architecture, proporcionando uma solução modular, testável e de fácil manutenção.

### Requisitos Funcionais

1. **Cotação de Frete**: Endpoint que recebe informações de volumes e CEP para cotação de frete
2. **Integração Externa**: Consumo da API do Frete Rápido para obtenção de cotações reais
3. **Persistência**: Armazenamento das cotações em banco de dados
4. **Métricas**: Endpoint para consulta de métricas sobre as cotações realizadas

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
│   ├── controllers/         # Controllers HTTP
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
- **Documentação**: Swagger
- **Testes**: Testify, Go-SQLMock
- **Containerização**: Docker e Docker Compose

## Instalação e Configuração

### Pré-requisitos
- Docker e Docker Compose

1. Clone o repositório:
   ```bash
   git clone https://github.com/thalesmacedo1/freterapido-backend-api.git
   cd freterapido-backend-api
   ```

2. Inicie os containers:
   ```bash
   docker-compose up -d
   ```
   ou com o comando
   ```bash
   make run
   ```
   A aplicação roda na seguinte URL e porta:
   ```bash
   http://localhost:3000/
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

## Documentação Swagger

A API utiliza o Swagger para documentação interativa dos endpoints. A documentação pode ser acessada em:

```
http://localhost:3000/swagger/index.html
```

### Gerando a documentação Swagger:

Para atualizar a documentação após mudanças no código, execute:

```bash
swag init -g api/cmd/api/main.go -o docs
```
