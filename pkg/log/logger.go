package log

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// Module register logger in DI container
var Module = fx.Provide(NewLogger)

// LoggerParams .
type LoggerParams struct {
	fx.In

	Level string `name:"log_level" optional:"true"`
}

// NewLogger gives new predefined logger
func NewLogger(p LoggerParams) (*logrus.Entry, error) {
	logger := logrus.New()

	if p.Level != "" {
		level, err := logrus.ParseLevel(p.Level)
		if err != nil {
			return nil, err
		}

		logger.SetLevel(level)
	}

	return logrus.NewEntry(logger), nil
}
