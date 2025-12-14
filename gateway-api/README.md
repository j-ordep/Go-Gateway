# Go Gateway API

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org/)
[![Chi Framework](https://img.shields.io/badge/Chi-v5.2.1-green.svg)](https://github.com/go-chi/chi)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)]()

Uma API REST moderna desenvolvida em Go para gerenciamento de contas e gateway de pagamentos, seguindo princípios de Clean Architecture e Domain-Driven Design (DDD).

## Características

- **Clean Architecture**: Separação clara de responsabilidades entre camadas
- **Domain-Driven Design**: Modelagem rica do domínio
- **API REST**: Endpoints RESTful para gerenciamento de contas
- **Thread Safety**: Controle de concorrência com mutexes
- **Testes Unitários**: Cobertura de testes com mocks
- **Arquitetura Hexagonal**: Desacoplamento entre domínio e infraestrutura

## Estrutura do Projeto

```
.
├── cmd/
│   └── app/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── domain/                  # Regras de negócio e entidades
│   │   ├── account.go          # Entidade Account
│   │   ├── repository.go       # Interface do repositório
│   │   └── errors.go          # Errors customizados do domínio
│   ├── dto/                    # Data Transfer Objects
│   │   └── account.go         # DTOs para Account
│   ├── repository/             # Implementação da camada de dados
│   │   └── account_repository.go
│   ├── service/                # Lógica de aplicação
│   │   ├── account_service.go
│   │   └── uptadete_test.go   # Testes unitários
│   └── web/                   # Camada de apresentação
│       ├── handler/
│       │   └── account_handler.go
│       └── server/
│           └── server.go
├── go.mod
├── go.sum
└── README.md
```

## Tecnologias Utilizadas

- **[Go 1.24.2](https://golang.org/)** - Linguagem de programação
- **[Chi Router v5](https://github.com/go-chi/chi)** - HTTP router e middleware
- **[Google UUID](https://github.com/google/uuid)** - Geração de UUIDs
- **[Testify](https://github.com/stretchr/testify)** - Framework de testes e mocks
- **SQL Database** - Persistência de dados

## Arquitetura

O projeto segue os princípios de **Clean Architecture** com as seguintes camadas:

### Domain Layer (Domínio)
- **Entities**: `Account` - Entidade principal com regras de negócio
- **Repositories**: Interfaces para acesso a dados
- **Errors**: Definição de erros específicos do domínio

### Application Layer (Aplicação)
- **Services**: `AccountService` - Orquestração de casos de uso
- **DTOs**: Objetos de transferência de dados

### Infrastructure Layer (Infraestrutura)
- **Repositories**: Implementação concreta para banco de dados
- **Web**: Handlers HTTP e configuração do servidor

## Funcionalidades

### Account Management

- **Criar Conta**: Criação de nova conta com API Key única
- **Consultar Conta**: Busca por API Key ou ID
- **Atualizar Saldo**: Operação thread-safe para atualização de saldo
- **Geração Automática**: API Key e timestamps automáticos

### Segurança

- **API Key Authentication**: Autenticação via header `X-API-KEY`
- **Thread Safety**: Mutex para operações de saldo
- **Validation**: Validação de dados de entrada

## Princípios de Design

### Clean Architecture
- **Dependency Inversion**: Interfaces definem contratos
- **Separation of Concerns**: Cada camada tem uma responsabilidade
- **Independence**: Camadas independentes e testáveis

### Domain-Driven Design
- **Rich Domain Model**: Entidades com comportamento
- **Ubiquitous Language**: Linguagem consistente
- **Bounded Context**: Contexto bem definido

### SOLID Principles
- **Single Responsibility**: Cada classe tem uma responsabilidade
- **Open/Closed**: Aberto para extensão, fechado para modificação
- **Liskov Substitution**: Subtipos substituíveis
- **Interface Segregation**: Interfaces específicas
- **Dependency Inversion**: Dependa de abstrações

**Variáveis de ambiente**
```bash
HTTP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gateway
DB_SSL_MODE=disable

# Configurações do Kafka
# Endereço do broker Kafka (pode ser uma lista separada por vírgulas para múltiplos brokers)
KAFKA_BROKER=localhost:9092

# Tópico para envio de transações pendentes (alto valor) para análise
KAFKA_PENDING_TRANSACTIONS_TOPIC=pending_transactions

# Tópico para recebimento dos resultados das análises de transações
KAFKA_TRANSACTIONS_RESULT_TOPIC=transactions_result

# Identificador do grupo de consumidores Kafka
# Deve ser único para cada instância do gateway quando executando em cluster
KAFKA_CONSUMER_GROUP_ID=gateway-group
```

## Rodando kafka

Producer (dentro do contêiner Kafka, usando listener interno `kafka:29092`):

```
 docker-compose exec kafka kafka-console-producer --bootstrap-server kafka:29092 --topic transactions_result
```

Consumer (dentro do contêiner Kafka):

```
 docker-compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic transactions_result --group group-result
```

Se quiser ler desde o início do tópico (apenas se o tópico estiver limpo):

```
 docker-compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic transactions_result --group group-result --from-beginning
```

Alternativa pelo host (WSL/Windows), sem `exec` (listener `localhost:9092`):

```
 kafka-console-producer --bootstrap-server localhost:9092 --topic transactions_result
 kafka-console-consumer --bootstrap-server localhost:9092 --topic transactions_result --group group-result
```

Exemplo de mensagem (uma única linha, sem prefixos `>` ou `#`):

```
{"invoice_id":"8561acea-bb75-48ee-aa0c-4e5b642de11b","status":"approved"}
```