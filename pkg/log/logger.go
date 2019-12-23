package log

import (
	"context"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/sirupsen/logrus"
)

const logPayload = true

// NewLogger gives new predefined logger
func NewLogger() *logrus.Entry {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel)

	return logrus.NewEntry(logger)
}

// LoggingPayloadDecider decide is need to log payload
func LoggingPayloadDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool { return logPayload }
}
