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

## Começando

### Pré-requisitos

- Go 1.24.2 ou superior
- Banco de dados SQL (PostgreSQL, MySQL, etc.)

### Instalação

1. **Clone o repositório**
```bash
git clone https://github.com/j-ordep/Go-Gateway.git
cd GO-GATEWAY-API
```

2. **Instale as dependências**
```bash
go mod download
```

3. **Configure o banco de dados**
```bash
# Configure suas variáveis de ambiente para conexão com BD
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=gateway
export DB_USER=your_user
export DB_PASSWORD=your_password
```

4. **Execute a aplicação**
```bash
go run cmd/app/main.go
```

## API Endpoints

### Accounts

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|--------------|
| `POST` | `/accounts` | Criar nova conta | Não |
| `GET` | `/accounts` | Buscar conta | API Key |

### Exemplos de Uso

#### Criar Conta
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@email.com"
  }'
```

**Resposta:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao@email.com",
  "api_key": "a1b2c3d4e5f6g7h8",
  "balance": 0,
  "created_at": "2025-08-05T10:30:00Z",
  "updated_at": "2025-08-05T10:30:00Z"
}
```

#### Buscar Conta
```bash
curl -X GET http://localhost:8080/accounts \
  -H "X-API-KEY: a1b2c3d4e5f6g7h8"
```

## Testes

### Executar Testes
```bash
# Todos os testes
go test ./...

# Testes com cobertura
go test -cover ./...

# Testes específicos
go test ./internal/service/
```

### Estrutura de Testes
- **Unit Tests**: Testes com mocks para isolamento
- **Integration Tests**: Testes de fluxo completo
- **Test Coverage**: Cobertura de código

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

## Roadmap

- [ ] **Authentication System**: Sistema de autenticação JWT
- [ ] **Payment Processing**: Integração com gateways de pagamento
- [ ] **Webhook Support**: Suporte a webhooks
- [ ] **Rate Limiting**: Limitação de requisições
- [ ] **Logging**: Sistema de logs estruturados
- [ ] **Monitoring**: Métricas e observabilidade
- [ ] **Docker Support**: Containerização
- [ ] **CI/CD Pipeline**: Automação de deploy

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## License

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## Autor

**João Pedro** - [@j-ordep](https://github.com/j-ordep)

---

**Star este projeto se ele foi útil para você!**