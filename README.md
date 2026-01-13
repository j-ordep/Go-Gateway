# Go Gateway API

Uma API REST moderna desenvolvida em Go para gerenciamento de contas e gateway de pagamentos, seguindo princípios de Clean Architecture e Domain-Driven Design (DDD).

## Características

- **Clean Architecture**: Separação clara de responsabilidades entre camadas
- **Domain-Driven Design**: Modelagem rica do domínio
- **API REST**: Endpoints RESTful para gerenciamento de contas
- **Thread Safety**: Controle de concorrência com mutexes
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
# Go Gateway API

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org/)
[![Chi Framework](https://img.shields.io/badge/Chi-v5.2.1-green.svg)](https://github.com/go-chi/chi)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)]()

API REST moderna escrita em Go para gerenciar contas, faturas e disparar análises antifraude. Este repositório é a raiz do monorepo do Gateway de Pagamentos construído durante a Imersão Full Stack & Full Cycle.

## Visão Geral do Monorepo

- **Gateway API (Go)**: serviço principal exposto ao frontend, aplica regras de negócio, grava dados em PostgreSQL e integra com Apache Kafka.
- **Microserviço Anti-Fraude (NestJS)**: consome transações suspeitas, processa a heurística de risco e devolve o veredito através de Kafka.
- **Frontend (Next.js)**: painel web que autentica via API Key, cria faturas e acompanha o status atualizado em tempo real.

## Arquitetura e Fluxo de Pagamentos

```
Next.js UI ──> Go Gateway (REST) ──> PostgreSQL
                         │
                         ├─(valor <= 10k)→ aprova/rejeita localmente
                         └─(valor > 10k)→ Kafka topic pending_transactions
                                             │
                                  NestJS Anti-Fraud Service
                                             │
                               Kafka topic transactions_result
                                             │
                             Go Gateway atualiza fatura e saldo
```

1. O frontend coleta a API Key e envia requisições ao gateway.
2. O serviço Go cria a fatura e roda `invoice.Process()`; valores até R$ 10.000 recebem decisão local probabilística (70% de aprovação) e já atualizam o saldo da conta ([gateway-api/internal/service/invoice_service.go](gateway-api/internal/service/invoice_service.go)).
3. Faturas acima de R$ 10.000 ficam em `pending` e disparam um evento `PendingTransaction` em `pending_transactions` via `KafkaProducer` ([gateway-api/internal/service/kafka.go](gateway-api/internal/service/kafka.go)).
4. O microserviço NestJS consome o evento em `nestjs-anti-fraud/src/invoices/invoices.consumer.ts` e aciona o `FraudService`, que roda as regras e publica `transactions_result` usando o cliente Confluent Kafka ([nestjs-anti-fraud/src/kafka/confluent-kafka-context.ts](nestjs-anti-fraud/src/kafka/confluent-kafka-context.ts)).
5. O `KafkaConsumer` Go escuta `transactions_result`, converte o payload para `TransactionResult` e atualiza o status/saldo da fatura via `ProcessTransactionResult`.

## Stack por Componente

| Serviço | Tecnologias | Pasta |
| --- | --- | --- |
| Gateway API | Go 1.24, Chi Router, PostgreSQL (database/sql), Segmentio Kafka, Testify | [gateway-api](gateway-api/README.md) |
| Anti-Fraude | NestJS 11, Prisma ORM, Confluent Kafka JS, Docker, TypeScript | [nestjs-anti-fraud](nestjs-anti-fraud/README.md) |
| Frontend | Next.js 15 App Router, Tailwind CSS, Shadcn UI, TypeScript | [next-frontend](next-frontend/README.md) |

## Componentes

### Gateway API (Go)

#### Principais características

- Clean Architecture + DDD com camadas `domain`, `service`, `repository` e `web` desacopladas.
- Autenticação por API Key (`X-API-KEY`) com mutexes garantindo operações de saldo thread-safe.
- Serviço de faturas capaz de aprovar/recusar localmente ou enviar eventos Kafka para análise externa.
- Handlers HTTP expostos via Chi Router; servidor configurado em [gateway-api/internal/web/server/server.go](gateway-api/internal/web/server/server.go).

#### Estrutura

```
gateway-api/
├── cmd/app/main.go          # bootstrap do servidor, DB e produtores/consumidores Kafka
├── internal/domain          # entidades, erros e eventos (Account, Invoice, Status, PendingTransaction)
├── internal/dto             # DTOs de entrada/saída
├── internal/repository      # implementações SQL
├── internal/service         # casos de uso (AccountService, InvoiceService, Kafka)
└── internal/web             # handlers REST e servidor HTTP
```

#### Tecnologias-chave

- Chi Router para rotas REST.
- `database/sql` + driver `pq` para PostgreSQL.
- Segmentio `kafka-go` como produtor/consumidor.
- `godotenv` para configuração local.
- Testify para testes unitários (ver `internal/service/*_test.go`).

#### Integração com Kafka

- `KafkaProducer` e `KafkaConsumer` compartilham `KafkaConfig` reutilizável com suporte a múltiplos brokers.
- Eventos:
  - `pending_transactions`: payload `PendingTransaction {account_id, invoice_id, amount}`.
  - `transactions_result`: payload `TransactionResult {invoice_id, status}` convertido para `domain.Status`.
- Consumidor roda em goroutine dedicada no `main`, aplicando `ProcessTransactionResult` e atualizando o saldo quando o status final é `approved`.

#### Variáveis de ambiente

```bash
HTTP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gateway
DB_SSL_MODE=disable

KAFKA_BROKER=localhost:9092                 # aceita lista separada por vírgula
KAFKA_PRODUCER_TOPIC=pending_transactions   # usado ao criar faturas
KAFKA_TRANSACTIONS_RESULT_TOPIC=transactions_result
KAFKA_CONSUMER_GROUP_ID=gateway-group
```

Use os utilitários Kafka descritos abaixo para inspecionar mensagens.

### Microserviço Anti-Fraude (NestJS)

- Código em [nestjs-anti-fraud](nestjs-anti-fraud/README.md).
- Expõe endpoints REST e um consumidor Kafka (`@EventPattern('pending_transactions')`).
- Prisma ORM conecta-se a um banco PostgreSQL isolado.
- Pode rodar só o consumidor com `npm run start:dev -- --entryFile cmd/kafka.cmd` dentro do container.
- Implementa regras de valor, frequência e comportamento de cartão para marcar `approved`, `rejected` ou manter `pending`.

### Frontend (Next.js)

- Código em [next-frontend](next-frontend/README.md).
- App Router + Tailwind + Shadcn UI com 4 telas (Login, Listagem, Detalhe, Criação de fatura).
- Integra com o gateway via API Key, invalida cache com revalidação e colore status (verde/amarelo/vermelho).
- Requer o serviço Go rodando primeiro, pois compartilha a rede Docker do compose principal.

## Operando Kafka

Producer (no container Kafka, listener `kafka:29092`):

```
docker-compose exec kafka kafka-console-producer --bootstrap-server kafka:29092 --topic transactions_result
```

Consumer (no container Kafka):

```
docker-compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic transactions_result --group group-result
```

Consumer lendo do início:

```
docker-compose exec kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic transactions_result --group group-result --from-beginning
```

Executando pelo host (listener `localhost:9092`):

```
kafka-console-producer --bootstrap-server localhost:9092 --topic transactions_result
kafka-console-consumer --bootstrap-server localhost:9092 --topic transactions_result --group group-result
```

Mensagem de exemplo:

```
{"invoice_id":"8561acea-bb75-48ee-aa0c-4e5b642de11b","status":"approved"}
```


- **Liskov Substitution**: Subtipos substituíveis
- **Interface Segregation**: Interfaces específicas
- **Dependency Inversion**: Dependa de abstrações
