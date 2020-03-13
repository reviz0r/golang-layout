package config

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// Module gives config params from file, env and flags
var Module = fx.Provide(NewConfig)

// NewConfig gives new predefined config provider
func NewConfig() (*viper.Viper, error) {
	config := viper.New()
	config.SetEnvKeyReplacer(strings.NewReplacer("__", "."))

	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("./configs")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	return config, nil
}
