package server

import (
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewGrpcServer gives new predefined grpc server
func NewGrpcServer(logger *logrus.Entry) *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcLogrus.StreamServerInterceptor(logger),
			grpcLogrus.PayloadStreamServerInterceptor(logger, log.LoggingPayloadDecider()),
			grpcRecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcLogrus.UnaryServerInterceptor(logger),
			grpcLogrus.PayloadUnaryServerInterceptor(logger, log.LoggingPayloadDecider()),
			grpcRecovery.UnaryServerInterceptor(),
		)),
	)
}
