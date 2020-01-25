package server

import (
	"context"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"go.uber.org/fx"
)

var GrpcLoggingPayloadModule = fx.Provide(LoggingPayloadDecider)

type GrpcLoggingParams struct {
	fx.In

	LogPayload bool `name:"grpc_log_payload" optional:"true"`
}

// LoggingPayloadDecider decide is need to log payload
func LoggingPayloadDecider(p GrpcLoggingParams) grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return p.LogPayload
	}
}
