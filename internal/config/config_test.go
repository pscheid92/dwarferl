package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGatherConfig(t *testing.T) {
	t.Run("successfully parses with defaults", func(t *testing.T) {
		_, err := GatherConfig()
		assert.NoErrorf(t, err, "unexpected error: %v", err)
	})

	t.Run("successfully reads database config", func(t *testing.T) {
		url := "this_is_a_test_url"
		pwd := "this_is_a_test_password"

		err := os.Setenv("DATABASE_URL", url)
		assert.NoErrorf(t, err, "unexpected error: %v", err)

		err = os.Setenv("DATABASE_PASSWORD", pwd)
		assert.NoErrorf(t, err, "unexpected error: %v", err)

		config, err := GatherConfig()
		assert.NoErrorf(t, err, "unexpected error: %v", err)
		assert.Equalf(t, url, config.DatabaseURL, "expected database url to be '%s', got %v", url, config.DatabaseURL)
		assert.Equal(t, pwd, config.DatabasePassword, "expected database password to be '%s', got %v", pwd, config.DatabasePassword)
	})

	t.Run("successfully read individualised env vars", func(t *testing.T) {
		err := os.Setenv("DWARFERL_USER", "this_is_a_test")
		assert.NoErrorf(t, err, "unexpected error: %v", err)

		config, err := GatherConfig()
		assert.NoErrorf(t, err, "unexpected error: %v", err)
		assert.Equal(t, "this_is_a_test", config.BasicAuthUser)
	})

	t.Run("successfully appends trailing slash to forwarded prefix", func(t *testing.T) {
		err := os.Setenv("FORWARDED_PREFIX", "/dummy")
		assert.NoErrorf(t, err, "unexpected error: %v", err)

		config, err := GatherConfig()
		assert.NoErrorf(t, err, "unexpected error: %v", err)
		assert.Equal(t, "/dummy/", config.ForwardedPrefix)
	})

	t.Run("fails if forwarded prefix does not start with a slash", func(t *testing.T) {
		err := os.Setenv("FORWARDED_PREFIX", "dummy")
		assert.NoErrorf(t, err, "unexpected error: %v", err)

		_, err = GatherConfig()
		assert.Errorf(t, err, "unexpected error: %v", err)
	})
}
