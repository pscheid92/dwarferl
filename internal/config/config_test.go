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
