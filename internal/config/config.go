package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	BasicAuthUser   string `mapstructure:"basic_auth_user"`
	BasicAuthSecret string `mapstructure:"basic_auth_secret"`
	ForwardedPrefix string `mapstructure:"forwarded_prefix"`
	SessionSecret   string `mapstructure:"session_secret"`
}

func GatherConfig() (Configuration, error) {
	// basic auth
	viper.SetDefault("basic_auth_user", "admin")
	viper.SetDefault("basic_auth_secret", "admin")

	// forwarded prefix
	viper.SetDefault("forwarded_prefix", "/")

	// session secret
	viper.SetDefault("session_secret", "secret")

	// environment variable bindings
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if !strings.HasPrefix(config.ForwardedPrefix, "/") {
		return Configuration{}, errors.New("forwarded_prefix must start with /")
	}

	if !strings.HasSuffix(config.ForwardedPrefix, "/") {
		config.ForwardedPrefix += "/"
	}

	return config, nil
}
