// internal/metrics/registry.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

type Registry struct {
    registry *prometheus.Registry
}

func NewRegistry() *Registry {
    return &Registry{
        registry: prometheus.NewRegistry(),
    }
}

func (r *Registry) Register(collector prometheus.Collector) error {
    return r.registry.Register(collector)
}

func (r *Registry) Handler() http.Handler {
    return promhttp.HandlerFor(r.registry, promhttp.HandlerOpts{})
}