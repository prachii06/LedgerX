package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// EventsConsumedTotal tracks the total number of Kafka events successfully consumed and persisted.
	EventsConsumedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ledgerx_events_consumed_total",
			Help: "Total number of Kafka events successfully consumed and persisted.",
		},
	)

	// EventsFailedTotal tracks the total number of Kafka events whose processing failed.
	EventsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ledgerx_events_failed_total",
			Help: "Total number of Kafka events whose processing failed.",
		},
	)

	// EventsDLQTotal tracks the total number of Kafka events successfully sent to the DLQ.
	EventsDLQTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ledgerx_events_dlq_total",
			Help: "Total number of Kafka events successfully sent to the dead letter queue.",
		},
	)

	// ReconciliationsTotal tracks reconciliation attempts by status.
	ReconciliationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ledgerx_reconciliations_total",
			Help: "Total number of reconciliations by status.",
		},
		[]string{"status"},
	)

	// EventProcessingDuration tracks the time taken to process a Kafka event.
	EventProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ledgerx_event_processing_duration_seconds",
			Help:    "Time taken to process a Kafka event.",
			Buckets: prometheus.DefBuckets,
		},
	)

	initOnce sync.Once
)

// Init registers all custom metrics with Prometheus.
// It uses sync.Once to ensure metrics are only registered once.
func Init() {
	initOnce.Do(func() {
		prometheus.MustRegister(EventsConsumedTotal)
		prometheus.MustRegister(EventsFailedTotal)
		prometheus.MustRegister(EventsDLQTotal)
		prometheus.MustRegister(ReconciliationsTotal)
		prometheus.MustRegister(EventProcessingDuration)
	})
}
