package message_service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type metrics struct {
	invalidFN        *prometheus.CounterVec
	handlingDuration prometheus.Histogram
}

func newMetrics() *metrics {
	return &metrics{
		invalidFN: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "fn_enricher",
				Subsystem: "fn_processor",
				Name:      "invalid_fn_count",
				Help:      "fn validation errors count",
			},
			[]string{"reason"},
		),
		handlingDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: "fn_enricher",
				Subsystem: "fn_processor",
				Name:      "fn_handling_duration",
				Help:      "fn handling duration histogram",
			}),
	}
}

func (m *metrics) incInvalidFN(err error) {
	m.invalidFN.WithLabelValues(err.Error()).Inc()
}

func (m *metrics) observe(t time.Duration) {
	m.handlingDuration.Observe(t.Seconds())
}
