package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/prachii06/LedgerX/internal/metrics"
	"github.com/prachii06/LedgerX/internal/models"
)

type EventHandler interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}

type DLQPublisher interface {
	PublishToDLQ(
		ctx context.Context,
		event *models.Event,
	) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler
	dlq     DLQPublisher
}

func NewConsumer(
	brokers string,
	groupID string,
	handler EventHandler,
	dlq DLQPublisher,
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
		dlq:     dlq,
	}
}

func (c *Consumer) Start(ctx context.Context) {

	log.Println("Kafka consumer started")

	for {

		// FetchMessage does NOT automatically commit the offset.
		message, err := c.reader.FetchMessage(ctx)

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

		if shouldCommit := c.processMessage(ctx, message); shouldCommit {
			if err := c.reader.CommitMessages(ctx, message); err != nil {
				log.Printf("Failed to commit Kafka message: %v", err)
			}
		}
	}
}

// processMessage processes a single Kafka message, handling retries and DLQ.
// It returns true if the message should be committed, false otherwise.
func (c *Consumer) processMessage(ctx context.Context, message kafka.Message) bool {
	start := time.Now()
	defer func() {
		metrics.EventProcessingDuration.Observe(time.Since(start).Seconds())
	}()

	var event models.Event

	if err := json.Unmarshal(message.Value, &event); err != nil {
		log.Printf("Failed to decode Kafka event: %v", err)
		// We cannot process an invalid event.
		// For now, leave the message uncommitted.
		return false
	}

	log.Printf(
		"Persisting event: transaction=%s type=%s sequence=%d",
		event.TransactionID,
		event.EventType,
		event.Sequence,
	)

	// Retry database persistence.
	const maxRetries = 3
	var persistErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		persistErr = c.handler.CreateEvent(ctx, &event)
		if persistErr == nil {
			break
		}

		log.Printf(
			"Failed to persist event %s (attempt %d/%d): %v",
			event.ID,
			attempt,
			maxRetries,
			persistErr,
		)

		if attempt < maxRetries {
			// Retry with increasing delay.
			select {
			case <-ctx.Done():
				return false
			case <-time.After(time.Duration(attempt) * time.Millisecond * 10): // Reduced delay for testing
			}
		}
	}

	// If all retries failed, send event to DLQ.
	if persistErr != nil {
		metrics.EventsFailedTotal.Inc()
		log.Printf("Event permanently failed after %d attempts: %s", maxRetries, event.ID)

		if err := c.dlq.PublishToDLQ(ctx, &event); err != nil {
			log.Printf("Failed to publish event %s to DLQ: %v", event.ID, err)
			// DO NOT commit the Kafka message. Kafka will redeliver it.
			return false
		}

		metrics.EventsDLQTotal.Inc()

		log.Printf("Event sent to DLQ: %s", event.ID)
		return true // Commit after successful DLQ publishing
	}

	metrics.EventsConsumedTotal.Inc()
	log.Printf("Event persisted successfully: transaction=%s type=%s", event.TransactionID, event.EventType)
	return true
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
