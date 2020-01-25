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

type ServerInterceptorParams struct {
	fx.In

	Logger                *logrus.Entry
	PayloadLoggingDecider grpcLogging.ServerPayloadLoggingDecider
}

type ServerInterceptorResult struct {
	fx.Out

	Option grpc.ServerOption `group:"grpc_server_options"`
}

func NewStreamServerInterceptors(p ServerInterceptorParams) ServerInterceptorResult {
	o := grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
		grpcLogrus.StreamServerInterceptor(p.Logger),
		grpcLogrus.PayloadStreamServerInterceptor(p.Logger, p.PayloadLoggingDecider),
		grpcRecovery.StreamServerInterceptor(),
		grpcValidator.StreamServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}

func NewUnaryServerInterceptors(p ServerInterceptorParams) ServerInterceptorResult {
	o := grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		grpcLogrus.UnaryServerInterceptor(p.Logger),
		grpcLogrus.PayloadUnaryServerInterceptor(p.Logger, p.PayloadLoggingDecider),
		grpcRecovery.UnaryServerInterceptor(),
		grpcValidator.UnaryServerInterceptor(),
	))

	return ServerInterceptorResult{Option: o}
}
