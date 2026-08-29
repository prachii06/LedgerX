package kafka

import (
	"context"
	"encoding/json"
	"github.com/prachii06/LedgerX/internal/models"
	"github.com/segmentio/kafka-go"
)

const EventTopic = "transaction-events"

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers),
			Topic:    EventTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishEvent(
	ctx context.Context,
	event *models.Event,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.TransactionID), //using TransactionID as the key for partitioning
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
