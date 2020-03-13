package logger

import (
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module register logger in DI container
var Module = fx.Provide(NewLogger, logrus.NewEntry)

// LoggerParams .
type LoggerParams struct {
	fx.In

	Output io.Writer `name:"logger_output" optional:"true"`
}

// NewLogger gives new predefined logger
func NewLogger(config *viper.Viper, p LoggerParams) (*logrus.Logger, error) {
	logger := logrus.New()

	if logFormatter := config.GetString("logger.formatter"); logFormatter != "" {
		switch logFormatter {
		case "text":
			logger.SetFormatter(new(logrus.TextFormatter))
		case "json":
			logger.SetFormatter(new(logrus.JSONFormatter))
		default:
			return nil, errors.New("unknown formatter")
		}
	}

	if logLevel := config.GetString("logger.level"); logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return nil, err
		}

		logger.SetLevel(level)
	}

	if config.GetBool("logger.no_lock") {
		logger.SetNoLock()
	}

	if p.Output != nil {
		logger.SetOutput(p.Output)
	}

	return logger, nil
}
