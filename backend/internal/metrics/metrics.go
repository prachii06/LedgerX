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

	// HTTPRequestsTotal tracks HTTP requests by method, route, and status code.
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ledgerx_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)

	// HTTPRequestDuration tracks HTTP request durations by method and route.
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "ledgerx_http_request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	// EventsProducedTotal tracks the total number of Kafka events successfully produced.
	EventsProducedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ledgerx_events_produced_total",
			Help: "Total number of Kafka events successfully published.",
		},
	)

	// ProducerErrorsTotal tracks Kafka producer errors.
	ProducerErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ledgerx_producer_errors_total",
			Help: "Total number of Kafka producer errors.",
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
		prometheus.MustRegister(HTTPRequestsTotal)
		prometheus.MustRegister(HTTPRequestDuration)
		prometheus.MustRegister(EventsProducedTotal)
		prometheus.MustRegister(ProducerErrorsTotal)
	})
}
