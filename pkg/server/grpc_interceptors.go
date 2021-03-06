package server

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpcOpenTracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
)

var InterceptorsModule = fx.Provide(NewStreamServerInterceptors, NewUnaryServerInterceptors)

type ServerInterceptorResult struct {
	fx.Out

	Option grpc.ServerOption `group:"grpc_server_options"`
}

func NewStreamServerInterceptors(logger *logrus.Entry, tracer opentracing.Tracer,
	PayloadLoggingDecider grpcLogging.ServerPayloadLoggingDecider) ServerInterceptorResult {
	o := grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
		grpcLogrus.StreamServerInterceptor(logger),
		grpcLogrus.PayloadStreamServerInterceptor(logger, PayloadLoggingDecider),
		grpcPrometheus.StreamServerInterceptor,
		grpcOpenTracing.OpenTracingStreamServerInterceptor(tracer),
		grpcRecovery.StreamServerInterceptor(),
		grpcValidator.StreamServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}

func NewUnaryServerInterceptors(logger *logrus.Entry, tracer opentracing.Tracer,
	PayloadLoggingDecider grpcLogging.ServerPayloadLoggingDecider) ServerInterceptorResult {
	o := grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		grpcLogrus.UnaryServerInterceptor(logger),
		grpcLogrus.PayloadUnaryServerInterceptor(logger, PayloadLoggingDecider),
		grpcPrometheus.UnaryServerInterceptor,
		grpcOpenTracing.OpenTracingServerInterceptor(tracer),
		grpcRecovery.UnaryServerInterceptor(),
		grpcValidator.UnaryServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}
