package internal

import (
	"github.com/spf13/viper"
	"strings"
)

type Configuration struct {
	BasicAuthUser   string `mapstructure:"basic_auth_user"`
	BasicAuthSecret string `mapstructure:"basic_auth_secret"`
	ForwardedPrefix string `mapstructure:"forwarded_prefix"`
	Database        struct {
		Host     string
		Port     uint16
		User     string
		Password string
		Name     string
	}
}

func GatherConfig() (Configuration, error) {
	// basic auth
	viper.SetDefault("basic_auth_user", "admin")
	viper.SetDefault("basic_auth_secret", "admin")

	// forwarded prefix
	viper.SetDefault("forward_prefix", "/")

	// database
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "postgres")

	// environment variable bindings
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("basic_auth_user", "DWARFERL_USER")
	viper.BindEnv("basic_auth_secret", "DWARFERL_SECRET")

	var config Configuration
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if !strings.HasSuffix(config.ForwardedPrefix, "/") {
		config.ForwardedPrefix += "/"
	}

	return config, nil
}
