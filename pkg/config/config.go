package config

import (
	"time"

	"go.uber.org/fx"
)

// Module gives config params from env
var Module = fx.Provide(NewConfig)

type Config struct {
	fx.Out

	// Database
	DatabaseDSN             string        `name:"database_dsn"`
	DatabaseConnMaxLifetime time.Duration `name:"database_conn_max_lifetime"`
	DatabaseMaxIdleConns    int           `name:"database_max_idle_conns"`
	DatabaseMaxOpenConns    int           `name:"database_max_open_conns"`
	DatabasePingOnStart     bool          `name:"database_ping_on_start"`

	// Logger
	LogFormatter string `name:"log_formatter"`
	LogLevel     string `name:"log_level"`
	LogNoLock    bool   `name:"log_no_lock"`
	LogPayload   bool   `name:"grpc_log_payload"`

	// other
	GrpcNetwork string `name:"grpc_network"`
	GrpcAddress string `name:"grpc_address"`

	HTTPNetwork string `name:"http_network"`
	HTTPAddress string `name:"http_address"`

	GatewayEndpoint     string `name:"gateway_user_service_endpoint"`
	GatewayEnumsAsInts  bool   `name:"gateway_marshaller_enums_as_ints"`
	GatewayEmitDefaults bool   `name:"gateway_marshaller_emit_defaults"`
	GatewayIndent       string `name:"gateway_marshaller_indent"`
	GatewayOrigName     bool   `name:"gateway_marshaller_orig_name"`
}

func NewConfig() Config {
	// TODO: use github.com/kelseyhightower/envconfig
	return Config{
		DatabaseDSN:         "user=postgres password=postgres host=localhost port=5432 sslmode=disable database=golang-layout",
		DatabasePingOnStart: true,

		LogFormatter: "text",
		LogLevel:     "trace",
		LogPayload:   true,

		GrpcNetwork: "tcp",
		GrpcAddress: ":50051",

		HTTPNetwork: "tcp",
		HTTPAddress: ":8081",

		GatewayEndpoint:     "localhost:50051",
		GatewayEmitDefaults: true,
	}
}
