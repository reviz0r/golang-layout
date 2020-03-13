package server

import (
	"context"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var GrpcLoggingPayloadModule = fx.Provide(LoggingPayloadDecider)

// LoggingPayloadDecider decide is need to log payload
func LoggingPayloadDecider(config *viper.Viper) grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return config.GetBool("logger.log_grpc_payload")
	}
}
