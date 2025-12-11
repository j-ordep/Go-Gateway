# Lógica / Regra de negócio 

# Lógica / Regra de negócio

## Visão Geral

Este projeto é um **gateway de pagamentos**, responsável por intermediar transações financeiras entre lojas virtuais (Accounts) e instituições financeiras.

---

## Variáveis de Ambiente

O sistema utiliza as seguintes variáveis de ambiente (ver arquivo `.env` na raiz do projeto):

```
HTTP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gateway
DB_SSL_MODE=disable

# Configurações do Kafka
KAFKA_BROKER=localhost:9092
KAFKA_PENDING_TRANSACTIONS_TOPIC=pending_transactions
KAFKA_TRANSACTIONS_RESULT_TOPIC=transactions_result
KAFKA_CONSUMER_GROUP_ID=gateway-group
```

## Entidades Principais

### Account
- Representa o cliente do gateway (ex: loja/empresa).
- Campos: `ID`, `Name`, `Email`, `APIKey`, `Balance`, `CreatedAt`, `UpdatedAt`.
- Cada Account possui um saldo e uma chave de API para autenticação.

### Invoice
- Representa uma cobrança/fatura gerada por uma Account para um pagamento específico.
- Campos: `ID`, `AccountId`, `Amount`, `Status`, `Description`, `PaymentType`, `CardLastDigits`, `CreatedAt`, `UpdatedAt`.
- Cada pagamento realizado gera uma Invoice.

### CreditCard
- Estrutura auxiliar para processar pagamentos via cartão de crédito.

---

## Fluxo de Criação de Pagamento

1. **Recebimento da requisição:**  
   O sistema recebe um pedido para criar uma Invoice, contendo dados do pagamento e da Account (APIKey).

2. **Identificação da Account:**  
   Busca a Account correspondente à APIKey fornecida.

3. **Criação da Invoice:**  
   - Valida o valor (não pode ser <= 0).
   - Salva os últimos 4 dígitos do cartão.
   - Define o status inicial como `pending`.

4. **Processamento da Invoice:**  
   - Se o valor for maior que 10.000, a Invoice permanece como `pending` (simulando análise manual/antifraude).
   - Se for menor ou igual a 10.000, o sistema sorteia (aleatoriamente) se a Invoice será `approved` (aprovada) ou `rejected` (rejeitada), com 70% de chance de aprovação.
   - O status é atualizado conforme o resultado.

5. **Atualização do saldo:**  
   Se a Invoice for aprovada, o saldo da Account é atualizado (adiciona o valor da Invoice).

6. **Persistência:**  
   A Invoice é salva no banco de dados.

---

## Observações

- **Status como tipo:**  
  O uso de um tipo customizado para Status (`type Status string`) garante segurança de tipo e centralização dos valores válidos.

- **APIKey:**  
  Usada para autenticar e identificar a Account que está criando a Invoice.

- **Repository e Service:**  
  - **Repository:** Responsável por interagir com o banco de dados.
  - **Service:** Implementa as regras de negócio, orquestrando as operações entre entidades e repositórios.

---

## Resumo

- **Account** = Loja/cliente do gateway.
- **Invoice** = Fatura/cobrança gerada para um pagamento.
- **Processamento** = Simula aprovação automática ou pendência para análise.
- **Saldo** = Atualizado quando uma cobrança é aprovada.