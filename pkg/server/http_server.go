package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var HTTPModule = fx.Provide(NewServeMux)

func NewServeMux(lc fx.Lifecycle, config *viper.Viper, logger *logrus.Entry) *http.ServeMux {
	mux := http.NewServeMux()

	address := config.GetString("http.address")
	s := http.Server{Addr: address, Handler: mux}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen(config.GetString("http.network"), address)
			if err != nil {
				return fmt.Errorf("cannot listen port %s %v", address, err)
			}

			go s.Serve(lis)
			logger.Debugf("http server started on port %s", address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := s.Shutdown(ctx)
			logger.Debug("http server is shutdown")
			return err
		},
	})

	return mux
}
