package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds Prometheus collectors used across the app.
type Metrics struct {
	RequestDuration *prometheus.HistogramVec
	RequestTotal    *prometheus.CounterVec
	OrdersCreated   prometheus.Counter
}

// New constructs and registers Prometheus metrics.
func New() *Metrics {
	m := &Metrics{
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "kabobfood",
			Name:      "request_latency_seconds",
			Help:      "HTTP request latency",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "path", "status"}),
		RequestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "kabobfood",
			Name:      "requests_total",
			Help:      "Total HTTP requests",
		}, []string{"method", "path", "status"}),
		OrdersCreated: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "kabobfood",
			Name:      "order_created_total",
			Help:      "Orders created",
		}),
	}

	prometheus.MustRegister(m.RequestDuration, m.RequestTotal, m.OrdersCreated)
	return m
}

// Middleware records HTTP metrics for Gin.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			m.RequestDuration.WithLabelValues(c.Request.Method, path, strconv.Itoa(c.Writer.Status())).Observe(v)
		}))
		c.Next()
		timer.ObserveDuration()
		status := strconv.Itoa(c.Writer.Status())
		m.RequestTotal.WithLabelValues(c.Request.Method, path, status).Inc()
	}
}
