package server

// http.Handle("/metrics", promhttp.Handler())

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

// PrometheusMetricsHandler .
var PrometheusMetricsHandler = fx.Invoke(HandlePrometheusMetrics)

// HandlePrometheusMetrics .
func HandlePrometheusMetrics(mux *http.ServeMux) {
	mux.Handle("/metrics", promhttp.Handler())
}
