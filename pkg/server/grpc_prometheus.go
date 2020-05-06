package server

import (
	"go.uber.org/fx"
	"google.golang.org/grpc"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
)

// PrometheusMetrics register grpc server in prometheus
var PrometheusMetrics = fx.Invoke(RegisterPrometheus)

// RegisterPrometheus .
func RegisterPrometheus(s *grpc.Server) {
	grpcPrometheus.Register(s)
}
