package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"

	"github.com/j-ordep/gateway/go-gateway/internal/domain/events"
	"github.com/segmentio/kafka-go"
)

// Anotação: por que existe `WithTopic`?
// - Objetivo: criar uma nova `KafkaConfig` mudando apenas o `Topic`, reaproveitando os mesmos `Brokers`.
// - Imutabilidade: evita mutar a config original; útil quando é compartilhada em múltiplos producers/consumers.
// - Conveniência: em cenários com um backend servindo vários fluxos (ou vários frontends),
//   todos usam o mesmo broker, mas publicam/consomem em tópicos diferentes (ex: "pending_transaction",
//   "transaction_result", "dead_letter"). `WithTopic` facilita compor essas variações.
// - Legibilidade: deixa explícito que só o tópico muda, reduzindo risco de esquecer campos ao copiar a struct.
// - Testes: em testes table-driven, gera configs variando apenas `Topic` sem impactar `Brokers`.
// - Exemplo de uso:
//     cfg := NewKafkaConfig()                         // Brokers definidos por env ou default
//     cfgPending := cfg.WithTopic("pending_transaction")
//     producerPending := NewKafkaProducer(cfgPending) // publica no tópico de pendências
//     cfgResult := cfg.WithTopic("transaction_result")
//     consumerResult := NewKafkaConsumer(cfgResult, "group-result", invoiceService)
//   Assim, reaproveitamos os mesmos brokers, variando somente o tópico para separar responsabilidades.

type KafkaProducerInterface interface {
	SendingPendingTransaction(ctx context.Context, event events.PendingTransaction) error
	Close() error
}

type KafkaConsumerInterface interface {
	Consume(ctx context.Context) error
	Close() error
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func (c *KafkaConfig) WithTopic(topic string) *KafkaConfig {
	return &KafkaConfig{
		Brokers: c.Brokers,
		Topic:   topic,
	}
}

func NewKafkaConfig() *KafkaConfig {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	topic := os.Getenv("KAFKA_PRODUCER_TOPIC")
	if topic == "" {
		topic = "pending_transactions"
	}

	return &KafkaConfig{
		Brokers: strings.Split(broker, ","),
		Topic:   topic,
	}
}

type KafkaProducer struct {
	writer  *kafka.Writer
	topic   string
	brokers []string
}

func NewKafkaProducer(config *KafkaConfig) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{}, //envia para a partição com menos bytes enfileirados. Otimiza throughput e reduz filas desbalanceadas sob carga variável.
	}
	// &kafka.Hash{}: escolhe partição com base no hash da Key. 
	// Garante afinidade de partição/ordem por chave (mensagens com mesma Key vão para a mesma partição).

	slog.Info("kafka producer iniciado", "brokers", config.Brokers, "topic", config.Topic)

	return &KafkaProducer{
		writer:  writer,
		topic:   config.Topic,
		brokers: config.Brokers,
	}
}

func (p *KafkaProducer) SendingPendingTransaction(ctx context.Context, event events.PendingTransaction) error {
	value, err := json.Marshal(event)
	if err != nil {
		slog.Error("erro ao converter evento para json", "error", err)
		return err
	}

	msg := kafka.Message{
		Value: value,
	}

	slog.Info("enviando mensagem para o kafka",
		"topic", p.topic,
		"message", string(value))

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		slog.Error("erro ao enviar mensagem para o kafka", "error", err)
		return err
	}

	slog.Info("mensagem enviada com sucesso para o kafka", "topic", p.topic)
	return nil
}

func (s *KafkaProducer) Close() error {
	slog.Info("fechando conexao com o kafka")
	return s.writer.Close()
}

type KafkaConsumer struct {
	reader         *kafka.Reader
	topic          string
	brokers        []string
	groupID        string
	invoiceService *InvoiceService
}

func NewKafkaConsumer(config *KafkaConfig, groupID string, invoiceService *InvoiceService) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		Topic:   config.Topic,
		GroupID: groupID,
	})

	slog.Info("kafka consumer iniciado",
		"brokers", config.Brokers,
		"topic", config.Topic,
		"group_id", groupID)

	return &KafkaConsumer{
		reader:         reader,
		topic:          config.Topic,
		brokers:        config.Brokers,
		groupID:        groupID,
		invoiceService: invoiceService,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("erro ao ler mensagem do kafka", "error", err)
			return err
		}

		var result events.TransactionResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			slog.Error("erro ao converter mensagem para TransactionResult", "error", err)
			continue
		}

		slog.Info("mensagem recebida do kafka",
			"topic", c.topic,
			"invoice_id", result.InvoiceID,
			"status", result.Status)

		// Processa o resultado da transação
		if err := c.invoiceService.ProcessTransactionResult(result.InvoiceID, result.ToDomainStatus()); err != nil {
			slog.Error("erro ao processar resultado da transação",
				"error", err,
				"invoice_id", result.InvoiceID,
				"status", result.Status)
			continue
		}

		slog.Info("transação processada com sucesso",
			"invoice_id", result.InvoiceID,
			"status", result.Status)
	}
}

func (c *KafkaConsumer) Close() error {
	slog.Info("fechando conexao com o kafka consumer")
	return c.reader.Close()
}
