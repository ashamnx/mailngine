package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mailngine_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mailngine_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	EmailsSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mailngine_emails_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{"org_id", "status"},
	)

	EmailsReceivedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "mailngine_emails_received_total",
			Help: "Total number of emails received",
		},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mailngine_active_connections",
			Help: "Number of active HTTP connections",
		},
	)
)
