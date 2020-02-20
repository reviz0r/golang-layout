package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/fx"
)

var HTTPModule = fx.Provide(NewServeMux)

type ServeMuxParams struct {
	fx.In

	Network string `name:"http_network"`
	Address string `name:"http_address"`
}

func NewServeMux(lc fx.Lifecycle, p ServeMuxParams) *http.ServeMux {
	mux := http.NewServeMux()

	s := http.Server{Addr: p.Address, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen(p.Network, p.Address)
			if err != nil {
				return fmt.Errorf("cannot listen port %s %v", p.Address, err)
			}

			go s.Serve(lis)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})

	return mux
}
