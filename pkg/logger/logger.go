package logger

import (
	"context"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module register logger in DI container
var Module = fx.Provide(NewLogger, logrus.NewEntry)

// NewLogger gives new predefined logger
func NewLogger(lc fx.Lifecycle, config *viper.Viper) (*logrus.Logger, error) {
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

	var outputFile *os.File
	if fileName := config.GetString("logger.output_file"); fileName != "" {
		var err error

		outputFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}

		logger.SetOutput(outputFile)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if outputFile != nil {
				return outputFile.Close()
			}

			return nil
		},
	})

	return logger, nil
}
