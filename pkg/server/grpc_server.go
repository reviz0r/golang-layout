package server

import (
	"context"
	"fmt"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/reviz0r/golang-layout/pkg/log"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module register grpc server in DI container
var Module = fx.Provide(NewGrpcServer)

// NewGrpcServer gives new predefined grpc server
func NewGrpcServer(lc fx.Lifecycle, logger *logrus.Entry) *grpc.Server {
	// TODO: in config
	port := ":50051"

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcLogrus.StreamServerInterceptor(logger),
			grpcLogrus.PayloadStreamServerInterceptor(logger, log.LoggingPayloadDecider()),
			grpcRecovery.StreamServerInterceptor(),
			grpcValidator.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcLogrus.UnaryServerInterceptor(logger),
			grpcLogrus.PayloadUnaryServerInterceptor(logger, log.LoggingPayloadDecider()),
			grpcRecovery.UnaryServerInterceptor(),
			grpcValidator.UnaryServerInterceptor(),
		)),
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", port)
			if err != nil {
				return fmt.Errorf("cannot listen grpc port %s %v", port, err)
			}

			go s.Serve(lis)
			logger.Debugf("grpc server started on port %s", port)
			return nil
		},

		OnStop: func(ctx context.Context) error {
			s.GracefulStop()
			logger.Debug("grpc server is shutdown")
			return nil
		},
	})

	return s
}
