package server

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module register grpc server in DI container
var Module = fx.Provide(NewGrpcServer)

// GrpcServerParams .
type GrpcServerParams struct {
	fx.In

	ServerOptions []grpc.ServerOption `group:"grpc_server_options"`
}

// NewGrpcServer gives new predefined grpc server
func NewGrpcServer(lc fx.Lifecycle, config *viper.Viper, logger *logrus.Entry, p GrpcServerParams) *grpc.Server {
	s := grpc.NewServer(p.ServerOptions...)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := config.GetString("grpc.address")
			lis, err := net.Listen(config.GetString("grpc.network"), address)
			if err != nil {
				return fmt.Errorf("cannot listen port %s %v", address, err)
			}

			go s.Serve(lis)
			logger.Debugf("grpc server started on port %s", address)
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
