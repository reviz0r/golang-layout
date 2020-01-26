package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/fx"
)

var HTTPModule = fx.Invoke(NewHTTPServer)

type HTTPServerParams struct {
	fx.In

	ProtoMux *runtime.ServeMux

	Network string `name:"http_network"`
	Address string `name:"http_address"`
}

func NewHTTPServer(lc fx.Lifecycle, p HTTPServerParams) {
	mux := http.NewServeMux()
	mux.Handle("/", p.ProtoMux)

	mux.HandleFunc("/docs/profile/swagger.json",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "api/proto/profile/profile.api.swagger.json")
		})

	s := http.Server{Addr: p.Address, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen(p.Network, p.Address)
			if err != nil {
				return fmt.Errorf("cannot listen port %s %v", p.Address, err)
			}

			go s.Serve(lis)
			// if err != nil && err != http.ErrServerClosed {
			// 	return fmt.Errorf("cannot serve http: %v", err)
			// }

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}
