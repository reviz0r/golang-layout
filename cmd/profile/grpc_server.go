package main

import (
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewGrpcServer gives new predefined grpc server
func NewGrpcServer(log *logrus.Entry) *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcLogrus.StreamServerInterceptor(log),
			grpcLogrus.PayloadStreamServerInterceptor(log, LoggingPayloadDecider()),
			grpcRecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcLogrus.UnaryServerInterceptor(log),
			grpcLogrus.PayloadUnaryServerInterceptor(log, LoggingPayloadDecider()),
			grpcRecovery.UnaryServerInterceptor(),
		)),
	)
}
