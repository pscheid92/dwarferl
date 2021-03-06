package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	ForwardedPrefix   string `mapstructure:"forwarded_prefix"`
	SessionSecret     string `mapstructure:"session_secret"`
	TemplatePath      string `mapstructure:"template_path"`
	AssetsPath        string `mapstructure:"assets_path"`
	GoogleClientKey   string `mapstructure:"google_client_key"`
	GoogleSecret      string `mapstructure:"google_secret"`
	GoogleCallbackURL string `mapstructure:"google_callback_url"`
}

func GatherConfig() (Configuration, error) {
	// forwarded prefix
	viper.SetDefault("forwarded_prefix", "/")

	// session secret
	viper.SetDefault("session_secret", "secret")

	// fs paths
	viper.SetDefault("template_path", "templates")
	viper.SetDefault("assets_path", "assets")

	// google login settings
	viper.SetDefault("google_client_key", "")
	viper.SetDefault("google_secret", "")
	viper.SetDefault("google_callback_url", "")

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
