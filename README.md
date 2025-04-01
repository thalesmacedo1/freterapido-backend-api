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
- [Contexto do Desafio](#contexto-do-desafio)
- [Arquitetura](#arquitetura)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação e Configuração](#instalação-e-configuração)
- [Rotas da API](#rotas-da-api)
- [Documentação Swagger](#documentação-swagger)
- [Testes](#testes)
- [Docker e Docker Compose](#docker-e-docker-compose)
- [Observações sobre a API do Frete Rápido](#observações-sobre-a-api-do-frete-rápido)
- [Conclusão e Trabalhos Futuros](#conclusão-e-trabalhos-futuros)

## Descrição

A API de Cotação de Frete é um sistema para consulta de valores de frete através de integrações com transportadoras. Desenvolvida com Go, segue os princípios de Domain-Driven Design (DDD) e Clean Architecture, proporcionando uma solução modular, testável e de fácil manutenção.

### Funcionalidades Principais

- Consulta de cotações de frete com diferentes transportadoras
- Métricas e análises estatísticas sobre as cotações realizadas
- Armazenamento de histórico de cotações para consulta posterior
- Integração com o serviço Frete Rápido para obtenção de cotações reais

## Contexto do Desafio

Este projeto foi desenvolvido como parte de um desafio técnico para demonstrar habilidades em desenvolvimento back-end com Go. O desafio consistia em criar uma API para cotação de fretes com os seguintes requisitos:

### Requisitos Funcionais

1. **Cotação de Frete**: Endpoint que recebe informações de volumes e CEP para cotação de frete
2. **Integração Externa**: Consumo da API do Frete Rápido para obtenção de cotações reais
3. **Persistência**: Armazenamento das cotações em banco de dados
4. **Métricas**: Endpoint para consulta de métricas sobre as cotações realizadas

### Requisitos Não-Funcionais

1. **Arquitetura**: Utilização de DDD e Clean Architecture
2. **Testes**: Implementação de testes unitários e de integração
3. **Documentação**: Documentação da API com Swagger
4. **Containerização**: Configuração com Docker e Docker Compose

### Solução Implementada

A solução consiste em uma API RESTful desenvolvida em Go com o framework Gin, seguindo princípios de Domain-Driven Design e Clean Architecture. A aplicação realiza consultas à API do Frete Rápido, armazena os resultados em um banco de dados PostgreSQL e oferece endpoints para cotação e análise de métricas.

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
- **Documentação**: Swagger
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
   go get -u github.com/swaggo/swag/cmd/swag
   go get -u github.com/swaggo/gin-swagger
   go get -u github.com/swaggo/files
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

## Documentação Swagger

A API utiliza o Swagger para documentação interativa dos endpoints. A documentação pode ser acessada em:

```
http://localhost:3000/swagger/index.html
```

### Recursos da documentação:

- **Visualização de endpoints**: Lista completa de endpoints com descrições
- **Modelos de dados**: Visualização de todas as estruturas de dados utilizadas
- **Testes interativos**: Possibilidade de testar os endpoints diretamente pela interface
- **Exemplos de resposta**: Exemplos de JSON para respostas bem-sucedidas

### Endpoints Documentados

A documentação Swagger inclui detalhes completos dos seguintes endpoints:

#### 1. POST /quote
- **Tags**: cotações
- **Descrição**: Retorna cotações de frete de diferentes transportadoras com base nos dados enviados
- **Parâmetros**: Corpo da requisição com dados do destinatário e volumes
- **Respostas**: 200 (sucesso), 400 (requisição inválida), 500 (erro interno)
- **Modelo de entrada**: domain.QuoteRequest
- **Modelo de saída**: domain.QuoteResponse

#### 2. GET /metrics
- **Tags**: métricas
- **Descrição**: Retorna métricas e estatísticas sobre as cotações de frete realizadas
- **Parâmetros**: Query parameter `last_quotes` (opcional) para limitar número de cotações
- **Respostas**: 200 (sucesso), 400 (parâmetro inválido), 500 (erro interno)
- **Modelo de saída**: domain.MetricsResponse

### Modelos de Dados

A documentação inclui definições completas dos seguintes modelos:

- **domain.QuoteRequest**: Solicitação para obter cotações de frete
- **domain.QuoteResponse**: Resposta com as cotações disponíveis
- **domain.Carrier**: Informações sobre uma transportadora específica
- **domain.Volume**: Detalhes de um volume para cotação
- **domain.MetricsResponse**: Resposta com métricas de cotações
- **domain.QuoteMetrics**: Métricas para uma transportadora específica
- **domain.CheapestAndMostExpensive**: Valores mínimos e máximos encontrados

### Implementação do Swagger

O Swagger foi implementado no projeto utilizando:

1. **gin-swagger**: Middleware para servir a documentação Swagger no Gin
2. **swaggo/swag**: Gerador de documentação Swagger a partir de anotações no código
3. **swaggo/files**: Arquivos estáticos para a UI do Swagger

A configuração do Swagger está presente em:
- `api/interfaces/routers/router.go`: Configuração da rota do Swagger
- `docs/swagger.json`: Arquivo de definição do Swagger gerado
- `docs/swagger.yaml`: Versão YAML da definição do Swagger

### Gerando a documentação Swagger:

Para atualizar a documentação após mudanças no código, execute:

```bash
swag init -g api/cmd/api/main.go -o docs
```

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

## Conclusão e Trabalhos Futuros

### Conclusão

Este projeto implementa uma API de Cotação de Frete seguindo boas práticas de desenvolvimento, arquitetura e documentação. As principais características do projeto incluem:

- **Design Robusto**: Arquitetura em camadas seguindo princípios de DDD e Clean Architecture
- **API Bem Documentada**: Documentação completa com Swagger
- **Testes**: Estrutura de testes unitários e de integração
- **Containerização**: Configuração com Docker para facilitar a execução e implantação
- **Integração Externa**: Conexão com serviço real de cotação de fretes

A API atende aos requisitos definidos no desafio, oferecendo endpoints para cotação de frete e análise de métricas, persistência em banco de dados e integração com o serviço Frete Rápido.

### Trabalhos Futuros

Os seguintes pontos foram identificados para melhorias e expansões futuras:

#### Aprimoramentos de Código
- **Completar testes**: Implementar mocks adequados para HTTP client e GORM
- **Corrigir testes de controllers**: Resolver problemas de compilação nos testes da camada de interface
- **Melhorar tratamento de erros**: Implementação mais robusta de tratamento de exceções

#### Novas Funcionalidades
- **Autenticação e Autorização**: Adicionar mecanismos de segurança como JWT
- **Rate Limiting**: Implementar controle de taxa para evitar sobrecarga da API
- **Cache**: Adicionar cache para cotações recentes para melhorar desempenho
- **Webhooks**: Implementar sistema de notificações para eventos específicos

#### Infraestrutura
- **CI/CD**: Configurar pipeline de integração e entrega contínua
- **Monitoramento**: Integrar com ferramentas como Prometheus e Grafana
- **Escalabilidade**: Preparar a aplicação para escalabilidade horizontal
- **Observabilidade**: Adicionar instrumentação para rastreamento de requisições

Este projeto serve como uma base sólida para uma aplicação de cotação de fretes e pode ser expandido em diversas direções conforme as necessidades do negócio evoluem.