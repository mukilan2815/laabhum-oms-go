package consumer

import (
    "context"
    "github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
    reader *kafka.Reader
}

func NewKafkaConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:  brokers,
        Topic:    topic,
        GroupID:  groupID,
        MinBytes: 10e3, // 10KB
        MaxBytes: 10e6, // 10MB
    })
    return &KafkaConsumer{reader: reader}
}

func (c *KafkaConsumer) StartProcessing(processMessage func(msg []byte) error) error {
    for {
        msg, err := c.reader.ReadMessage(context.Background())
        if err != nil {
            return err
        }
        if err := processMessage(msg.Value); err != nil {
            return err
        }
    }
}

func (c *KafkaConsumer) Close() error {
    return c.reader.Close()
}
