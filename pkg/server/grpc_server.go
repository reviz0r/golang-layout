package server

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// Module register grpc server in DI container
var Module = fx.Provide(NewGrpcServer)

// GrpcServerParams .
type GrpcServerParams struct {
	fx.In

	Logger *logrus.Entry

	Network string `name:"grpc_network"`
	Address string `name:"grpc_address"`

	ServerOptions []grpc.ServerOption `group:"grpc_server_options"`
}

// NewGrpcServer gives new predefined grpc server
func NewGrpcServer(lc fx.Lifecycle, p GrpcServerParams) *grpc.Server {
	s := grpc.NewServer(p.ServerOptions...)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen(p.Network, p.Address)
			if err != nil {
				return fmt.Errorf("cannot listen port %s %v", p.Address, err)
			}

			go s.Serve(lis)
			p.Logger.Debugf("grpc server started on port %s", p.Address)
			return nil
		},

		OnStop: func(ctx context.Context) error {
			s.GracefulStop()
			p.Logger.Debug("grpc server is shutdown")
			return nil
		},
	})

	return s
}
