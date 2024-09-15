package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(message []byte) error
	Close() error
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer using segmentio/kafka-go
func NewProducer(brokers []string, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{writer: writer}
}

// SendMessage publishes a message to the Kafka topic
func (p *KafkaProducer) SendMessage(message []byte) error {
	err := p.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// Shutdown closes the Kafka producer connection
func (p *KafkaProducer) Shutdown() error {
	return p.writer.Close()
}

// NewConsumer creates a new Kafka consumer using segmentio/kafka-go
func NewConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaConsumer{reader: reader}
}

// Consume reads messages from the Kafka topic
func (c *KafkaConsumer) Consume() ([]byte, error) {
	msg, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

// Close shuts down the Kafka consumer
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
