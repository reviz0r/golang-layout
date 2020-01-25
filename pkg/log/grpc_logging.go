package log

import (
	"context"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
)

const logPayload = true

// LoggingPayloadDecider decide is need to log payload
func LoggingPayloadDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return logPayload }
}
