package internal

import (
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	BasicAuthUser   string `mapstructure:"basic_auth_user"`
	BasicAuthSecret string `mapstructure:"basic_auth_secret"`
	ForwardedPrefix string `mapstructure:"forwarded_prefix"`
}

func GatherConfig() (Configuration, error) {
	viper.SetDefault("basic_auth_user", "admin")
	viper.SetDefault("basic_auth_secret", "admin")
	viper.SetDefault("forward_prefix", "/")

	viper.BindEnv("basic_auth_user", "DWARFERL_USER")
	viper.BindEnv("basic_auth_secret", "DWARFERL_SECRET")
	viper.BindEnv("forwarded_prefix", "FORWARDED_PREFIX")

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if !strings.HasSuffix(config.ForwardedPrefix, "/") {
		config.ForwardedPrefix += "/"
	}

	return config, nil
}
