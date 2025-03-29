# FreteRápido Backend Test

![Go](https://img.shields.io/badge/Go-1.20-blue)
![Gin](https://img.shields.io/badge/Gin-Framework-brightgreen)
![PostgreSQL](https://img.shields.io/badge/postgresql-4169e1?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Containerization-blue)

Este projeto é uma implementação de uma API RESTful para listagem e consulta de dados sobre a COVID-19.

## Sumário

- [Descrição](#descrição)
- [Problemas](#problemas)
- [Requisitos Funcionais](#requisitos-funcionais)
- [Arquitetura](#arquitetura)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Instalação e Configuração](#instalação-e-configuração)
- [Importação do CSV](#importação-do-csv)
- [Execução](#execução)
- [Rotas da API](#rotas-da-api)
- [Validação de Documentos (CPF/CNPJ)](#validação-de-documentos-cpfcnpj)
- [Testes e Cobertura](#testes-e-cobertura)
- [Docker e Docker Compose](#docker-e-docker-compose)
- [Melhorias Futuras](#melhorias-futuras)

---

## Descrição
A API de Dados COVID oferece acesso a informações relacionadas à vacinação da COVID-19. Com esta API, desenvolvedores podem recuperar dados sobre casos confirmados, óbitos, taxas de vacinação e outras métricas relevantes.

Principais funcionalidades:
Dados Globais e Regionais: Acesso a estatísticas de diversos países, estados e municípios.
Atualizações em Tempo Real: Informações constantemente atualizadas para refletir a situação atual da pandemia.
Histórico de Dados: Consulta de dados históricos para análises temporais e identificação de tendências.
Filtros por País, Região ou Data: Possibilidade de filtrar dados especificamente por país, região ou data desejada.

```
freterapido-backend-api/
│
├── domain/                # Camada de Domínio
│   ├── entities/          # Modelos de domínio puros
│   └── interfaces/        # Interfaces de repositórios e serviços DEFINIDAS NO DOMÍNIO
│       ├── repositories/  # Contratos (interfaces) de repositórios
│       └── services/      # Contratos (interfaces) de serviços de domínio
│
├── application/           # Camada de Infraestrutura
|   └── services/          # Serviços de aplicação
|    └── interfaces/       # Interfaces adicionais de aplicação
│
└── infrastructure/        # Camada de Aplicação
│   └── repositories/      
│       └── postgres/      
│       └── redis/         
```

---

## Requisitos Funcionais

1. Total Acumulado de Casos e Mortes por País e Data.
2. Número de Pessoas Vacinadas com Pelo Menos Uma Dose.
3. Vacinas Utilizadas e Data de Início de Aplicação.
4. País com o Maior Número de Casos Acumulados até uma Data.
5. Vacina Mais Utilizada em uma Região Específica.

---

## Arquitetura

A solução adota Clean Architecture e DDD, separando bem as camadas:

- **domain**: Entidades e regras de negócio puras.
- **usecase**: Casos de uso que orquestram a lógica de negócio.
- **infrastructure**: Implementações concretas de banco de dados, importação e validações.
- **interface**: Handlers HTTP e roteamento.
- **cmd/server**: Ponto de entrada da aplicação.

### Benefícios da Arquitetura Adotada
- **Separação de Preocupações**: Cada camada possui responsabilidades bem definidas, facilitando o entendimento e a manutenção do código.
- **Escalabilidade**: A modularidade permite que novas funcionalidades sejam adicionadas sem impactar negativamente as camadas existentes.
- **Testabilidade**: A independência entre as camadas facilita a realização de testes unitários e de integração, assegurando a qualidade do código.
- **Flexibilidade**: Utilizando interfaces na camada de domínio, é possível substituir ou atualizar implementações de infraestrutura com mínima alteração no restante do sistema.
- **Reutilização de Código**: Componentes bem definidos podem ser reutilizados em diferentes partes da aplicação ou até mesmo em outros projetos.
---

## Tecnologias Utilizadas

- **Linguagem**: Go (golang) 1.23+
- **Framework HTTP**: Gin
- **Banco de Dados**: PostgreSQL, Redis
- **Testes**: `testing` nativo do Go e a lib Testify
- **Docker**: Para containerização das duas aplicações e dos bancos de dados necessários

---

## Instalação e Configuração

1. **Pré-requisitos**:
   - Go 1.23+
   - Docker e Docker Compose
   - PostgreSQL

2. **Clonar o repositório**:
   ```bash
   git clone https://github.com/thalesmacedo1/freterapido-backend-api.git
   cd freterapido-backend-api
   ```

3. **Variáveis de Ambiente**:
   Os exemplos estão contidos no arquivo `.env.example` na raiz do projeto.

4. **Dependências**:
Rodando manualmente na pasta principal de cada uma das aplicações:
   ```bash
   go mod tidy
   ```

---

## Execução

### Sem Docker

- Certifique-se que o PostgreSQL está rodando e as variáveis de ambiente apontam para o banco de dados correto.
- Rode o servidor:
  ```bash
  go run cmd/server/main.go
  ```

A aplicação estará disponível em `http://localhost:8080`.

### Com Docker

Para executar com containers docker, rode:
   ```bash
   make run
   ```

---

## Rotas da API
Importante: o formato da data é YYYY-MM-DD

- **GET `/api/v1/countries/:countryCode/covid/:date`**  
Descrição: Obtém o total de casos confirmados e mortes em um país específico em uma data determinada.

  **Exemplo de resposta**:
  ```json
  {
    "error": "Failed to retrieve COVID totals."
  }
  ```


- **GET `/api/v1/countries/:countryCode/vaccinations/:date`**  
Descrição: Retorna a quantidade de indivíduos vacinados com pelo menos uma dose em um país específico em uma data determinada.

  **Exemplo de resposta**:
  ```json
  {
    "error": "Failed to retrieve vaccinated people data."
  }
  ```

- **GET `/api/v1/countries/:countryCode/vaccines`**  
Descrição: Lista as vacinas utilizadas em um país específico e as datas de início de aplicação.
Resultados para BRA (Brasil):

  **Exemplo de resposta**:
  ```json
  [
    {
      "vaccine": {
        "Product": "SII - Covishield",
        "Company": "Serum Institute of India",
        "Vaccine": "Covishield"
      },
      "start_date": "0001-01-01"
    },
    {
      "vaccine": {
        "Product": "Sinovac - CoronaVac",
        "Company": "Sinovac",
        "Vaccine": "Coronavac"
      },
      "start_date": "0001-01-01"
    },
    {
      "vaccine": {
        "Product": "Janssen - Ad26.COV 2-S",
        "Company": "Janssen Pharmaceuticals",
        "Vaccine": "Ad26.COV 2-S"
      },
      "start_date": "0001-01-01"
    },
    {
      "vaccine": {
        "Product": "Pfizer BioNTech - Comirnaty",
        "Company": "Pfizer BioNTech",
        "Vaccine": "Comirnaty"
      },
      "start_date": "0001-01-01"
    },
    {
      "vaccine": {
        "Product": "AstraZeneca - Vaxzevria",
        "Company": "AstraZeneca",
        "Vaccine": "Vaxzevria"
      },
      "start_date": "0001-01-01"
    }
  ]

  ```

- **GET `/api/v1/countries/highest-cases?:date`**  
Descrição: Identifica o país com o maior número de casos acumulados até uma data específica.

  **Exemplo de resposta**:
  ```json
  {
    "error": "Failed to get country with most cases"
  }
  ```

- **GET `/api/v1/regions/:regionName/vaccines/most-used`**  
Descrição: Retorna a vacina mais utilizada em uma determinada região.

  **Exemplo de resposta**:
  ```json
  {
    "error": "Failed to retrieve most used vaccine."
  }
  ```

---

## Testes e Cobertura

Para rodar os testes:

```bash
make test
```
---

## Docker e Docker Compose

O projeto inclui um `docker-compose.yml` e dois `Dockerfile`.

- O banco ficará disponível em `http://localhost:7474/`, com usuário `PostgreSQL` e senha padrão contida no `.env.example` 

A aplicação Go pode ser executada localmente conectando-se ao container do PostgreSQL conforme as variáveis de ambiente definidas no `.env.example`.

---

## Melhorias Futuras

- Implementar testes de integração para as rotas HTTP.
- Adicionar autenticação e autorização se preciso.
- Otimizar a importação, caso o CSV seja muito grande.
- Criar uma rota para importar o CSV sob demanda.