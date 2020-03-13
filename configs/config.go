package configs

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var Module = fx.Invoke(NewConfig)

func NewConfig(config *viper.Viper) {
	config.SetDefault("database.dsn", "host=localhost user=postgres sslmode=disable")
	config.SetDefault("grpc.network", "tcp")
	config.SetDefault("grpc.address", ":50051")
	config.SetDefault("http.network", "tcp")
	config.SetDefault("http.address", ":8081")
}
