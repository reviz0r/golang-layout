package log

import (
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// Module register logger in DI container
var Module = fx.Provide(NewLogger, logrus.NewEntry)

// LoggerParams .
type LoggerParams struct {
	fx.In

	Formatter string    `name:"log_formatter" optional:"true"`
	Level     string    `name:"log_level" optional:"true"`
	NoLock    bool      `name:"log_no_lock" optional:"true"`
	Output    io.Writer `name:"log_output" optional:"true"`
}

// NewLogger gives new predefined logger
func NewLogger(p LoggerParams) (*logrus.Logger, error) {
	logger := logrus.New()

	if p.Formatter != "" {
		switch p.Formatter {
		case "text":
			logger.SetFormatter(new(logrus.TextFormatter))
		case "json":
			logger.SetFormatter(new(logrus.JSONFormatter))
		default:
			return nil, errors.New("unknown formatter")
		}
	}

	if p.Level != "" {
		level, err := logrus.ParseLevel(p.Level)
		if err != nil {
			return nil, err
		}

		logger.SetLevel(level)
	}

	if p.NoLock {
		logger.SetNoLock()
	}

	if p.Output != nil {
		logger.SetOutput(p.Output)
	}

	return logger, nil
}
