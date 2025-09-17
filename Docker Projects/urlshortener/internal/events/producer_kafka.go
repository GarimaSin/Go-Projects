package events

import (
    "context"
    "encoding/json"
    "time"

    "github.com/segmentio/kafka-go"
)

type Producer interface {
    PublishClick(ctx context.Context, ev interface{}) error
    Close() error
}

type KafkaProducer struct {
    w *kafka.Writer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
    w := kafka.NewWriter(kafka.WriterConfig{
        Brokers:  brokers,
        Topic:    "clicks",
        Balancer: &kafka.Hash{},
        Async:    true,
    })
    return &KafkaProducer{w: w}, nil
}

func (p *KafkaProducer) PublishClick(ctx context.Context, ev interface{}) error {
    b, err := json.Marshal(ev)
    if err != nil {
        return err
    }
    return p.w.WriteMessages(ctx, kafka.Message{
        Key:   []byte(time.Now().Format(time.RFC3339Nano)),
        Value: b,
    })
}

func (p *KafkaProducer) Close() error {
    return p.w.Close()
}
