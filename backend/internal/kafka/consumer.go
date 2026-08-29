package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/prachii06/LedgerX/internal/models"
)

type EventHandler interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler
}

func NewConsumer(
	brokers string,
	groupID string,
	handler EventHandler,
) *Consumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokers},
		Topic:       EventTopic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {

	log.Println("Kafka consumer started")

	for {

		message, err := c.reader.ReadMessage(ctx)

		if err != nil {

			if ctx.Err() != nil {
				return
			}

			log.Printf(
				"Kafka consumer error: %v",
				err,
			)

			continue
		}

		log.Printf(
			"Kafka message received: topic=%s partition=%d offset=%d",
			message.Topic,
			message.Partition,
			message.Offset,
		)

		var event models.Event

		if err := json.Unmarshal(
			message.Value,
			&event,
		); err != nil {

			log.Printf(
				"Failed to decode Kafka event: %v",
				err,
			)

			continue
		}

		log.Printf(
			"Persisting event: transaction=%s type=%s sequence=%d",
			event.TransactionID,
			event.EventType,
			event.Sequence,
		)

		if err := c.handler.CreateEvent(
			ctx,
			&event,
		); err != nil {

			log.Printf(
				"Failed to persist event %s: %v",
				event.ID,
				err,
			)

			continue
		}

		log.Printf(
			"Event persisted successfully: transaction=%s type=%s",
			event.TransactionID,
			event.EventType,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}