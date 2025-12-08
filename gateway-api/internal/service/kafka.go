package service

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/j-ordep/gateway/go-gateway/internal/domain/events"
	"github.com/segmentio/kafka-go"
)

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

// WithTopic cria uma nova configuração com um tópico diferente
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
		topic = "pending_transaction"
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
		Balancer: &kafka.LeastBytes{},
	}

	slog.Info("kafka producer iniciado", "brokers", config.Brokers, "topic", config.Topic)

	return &KafkaProducer{
		writer:  writer,
		topic:   config.Topic,
		brokers: config.Brokers,
	}
}

func (p *KafkaProducer) SendingPendingTransaction(ctx context.Context, event events.PendingTransaction) error { return nil }

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

func (c *KafkaConsumer) Consume(ctx context.Context) error {return nil}

func (c *KafkaConsumer) Close() error {
	slog.Info("fechando conexao com o kafka consumer")
	return c.reader.Close()
}
