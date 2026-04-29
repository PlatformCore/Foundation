package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsConfig configures Prometheus metrics collection.
type MetricsConfig struct {
	// Registerer is the Prometheus registerer. Defaults to prometheus.DefaultRegisterer.
	Registerer prometheus.Registerer
	// Namespace is the metric name prefix. Default: "http".
	Namespace string
	// Subsystem is the metric name subsystem. Default: "server".
	Subsystem string
	// Buckets are the latency histogram buckets in seconds.
	Buckets []float64
	// LabelExtractors allows adding custom labels per-request.
	// Keys must match labels registered in ConstLabels or via LabelNames.
	LabelNames     []string
	LabelExtractor func(r *http.Request, status int) []string
	// SkipPaths lists paths excluded from metrics collection.
	SkipPaths []string
}

// serverMetrics holds all registered Prometheus metrics.
type serverMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
	responseSizeBytes *prometheus.HistogramVec
}

func newServerMetrics(cfg MetricsConfig) *serverMetrics {
	ns := cfg.Namespace
	if ns == "" {
		ns = "http"
	}
	sub := cfg.Subsystem
	if sub == "" {
		sub = "server"
	}
	buckets := cfg.Buckets
	if len(buckets) == 0 {
		buckets = prometheus.DefBuckets
	}

	baseLabels := []string{"method", "path", "status"}
	baseLabels = append(baseLabels, cfg.LabelNames...)

	reg := cfg.Registerer
	if reg == nil {
		reg = prometheus.DefaultRegisterer
	}

	m := &serverMetrics{}

	m.requestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "requests_total",
		Help:      "Total number of HTTP requests processed.",
	}, baseLabels)

	m.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency distribution.",
		Buckets:   buckets,
	}, baseLabels)

	m.requestsInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "requests_in_flight",
		Help:      "Number of HTTP requests currently being processed.",
	})

	m.responseSizeBytes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      "response_size_bytes",
		Help:      "HTTP response size distribution in bytes.",
		Buckets:   prometheus.ExponentialBuckets(100, 10, 7),
	}, baseLabels)

	reg.MustRegister(
		m.requestsTotal,
		m.requestDuration,
		m.requestsInFlight,
		m.responseSizeBytes,
	)

	return m
}

// WithMetrics returns a Prometheus metrics middleware.
// It tracks request count, latency, in-flight requests, and response size.
func WithMetrics(cfg MetricsConfig) Middleware {
	m := newServerMetrics(cfg)

	skipSet := make(map[string]struct{}, len(cfg.SkipPaths))
	for _, p := range cfg.SkipPaths {
		skipSet[p] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, skip := skipSet[r.URL.Path]; skip {
				next.ServeHTTP(w, r)
				return
			}

			m.requestsInFlight.Inc()
			defer m.requestsInFlight.Dec()

			start := time.Now()
			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			statusStr := strconv.Itoa(rw.status)

			baseLabels := []string{r.Method, r.URL.Path, statusStr}
			if cfg.LabelExtractor != nil {
				baseLabels = append(baseLabels, cfg.LabelExtractor(r, rw.status)...)
			} else {
				for range cfg.LabelNames {
					baseLabels = append(baseLabels, "")
				}
			}

			m.requestsTotal.WithLabelValues(baseLabels...).Inc()
			m.requestDuration.WithLabelValues(baseLabels...).Observe(duration)
			m.responseSizeBytes.WithLabelValues(baseLabels...).Observe(float64(rw.size))
		})
	}
}
