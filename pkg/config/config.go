package config

import "go.uber.org/fx"

// Module gives config params from env
var Module = fx.Provide(NewConfig)

type Config struct {
	fx.Out

	DatabaseDSN         string `name:"database_dsn"`
	DatabasePingOnStart bool   `name:"database_ping_on_start"`

	LogLevel   string `name:"log_level"`
	LogPayload bool   `name:"grpc_log_payload"`

	GrpcNetwork string `name:"grpc_network"`
	GrpcAddress string `name:"grpc_addtess"`
}

func NewConfig() Config {
	// TODO: use github.com/kelseyhightower/envconfig
	return Config{
		DatabaseDSN:         "user=postgres password=postgres host=localhost port=5432 sslmode=disable database=golang-layout",
		DatabasePingOnStart: true,

		LogLevel:   "trace",
		LogPayload: true,

		GrpcNetwork: "tcp",
		GrpcAddress: ":50051",
	}
}
