package kafka

import (
	"context"
	"encoding/json"

	"github.com/prachii06/LedgerX/internal/models"
	"github.com/segmentio/kafka-go"
)

const (
	EventTopic      = "transaction-events"
	DeadLetterTopic = "transaction-events-dlq"
)

type Producer struct {
	writer    *kafka.Writer
	dlqWriter *kafka.Writer
}

func NewProducer(brokers string) *Producer {

	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers),
			Topic:                  EventTopic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},

		dlqWriter: &kafka.Writer{
			Addr:                   kafka.TCP(brokers),
			Topic:                  DeadLetterTopic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
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
		Key:   []byte(event.TransactionID),
		Value: payload,
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *Producer) PublishToDLQ(
	ctx context.Context,
	event *models.Event,
) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(event.TransactionID),
		Value: payload,
	}

	return p.dlqWriter.WriteMessages(ctx, message)
}

func (p *Producer) Close() error {

	if err := p.writer.Close(); err != nil {
		return err
	}

	return p.dlqWriter.Close()
}
