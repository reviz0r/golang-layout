package mockserver

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Module register mock grpc server in DI container
var Module = fx.Provide(NewServer)

// NewServer gives new mocked grpc server
func NewServer(lc fx.Lifecycle) (*grpc.Server, *grpc.ClientConn, error) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	// create local grpc server
	s := grpc.NewServer()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// connect to local grpc server
	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.Serve(lis)
			return nil
		},

		OnStop: func(ctx context.Context) error {
			s.Stop()
			return conn.Close()
		},
	})

	return s, conn, nil
}
