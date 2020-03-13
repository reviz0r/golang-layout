package config

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// DefaultValues set dafault config values in DI container
var DefaultValues = fx.Invoke(SetConfigDefaults)

// SetConfigDefaults define default values for app config
func SetConfigDefaults(config *viper.Viper) {
	config.SetDefault("database.dsn", "host=localhost user=postgres sslmode=disable")
	config.SetDefault("grpc.network", "tcp")
	config.SetDefault("grpc.address", ":50051")
	config.SetDefault("http.network", "tcp")
	config.SetDefault("http.address", ":80")
}
