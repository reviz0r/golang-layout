package server

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
)

var InterceptorsModule = fx.Provide(NewStreamServerInterceptors, NewUnaryServerInterceptors)

type ServerInterceptorResult struct {
	fx.Out

	Option grpc.ServerOption `group:"grpc_server_options"`
}

func NewStreamServerInterceptors(logger *logrus.Entry,
	PayloadLoggingDecider grpcLogging.ServerPayloadLoggingDecider) ServerInterceptorResult {
	o := grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
		grpcLogrus.StreamServerInterceptor(logger),
		grpcLogrus.PayloadStreamServerInterceptor(logger, PayloadLoggingDecider),
		grpcRecovery.StreamServerInterceptor(),
		grpcValidator.StreamServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}

func NewUnaryServerInterceptors(logger *logrus.Entry,
	PayloadLoggingDecider grpcLogging.ServerPayloadLoggingDecider) ServerInterceptorResult {
	o := grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		grpcLogrus.UnaryServerInterceptor(logger),
		grpcLogrus.PayloadUnaryServerInterceptor(logger, PayloadLoggingDecider),
		grpcRecovery.UnaryServerInterceptor(),
		grpcValidator.UnaryServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}
