package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Prometheus struct {
	RequestDuration *prometheus.HistogramVec
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		RequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		}, []string{"path"}),
	}
}

func (p *Prometheus) Setup() {
	prometheus.MustRegister(p.RequestDuration)
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9090", nil)
}

func (p *Prometheus) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(p.RequestDuration.WithLabelValues(r.URL.Path))
		defer timer.ObserveDuration()
		next.ServeHTTP(w, r)
	})
}
