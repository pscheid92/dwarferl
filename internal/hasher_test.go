package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlHasher(t *testing.T) {
	tt := []struct {
		user  User
		url   string
		short string
	}{
		{
			DummyUser,
			"https://www.google.com",
			"hYhahA",
		},
		{
			DummyUser,
			"https://patrickscheid.de",
			"5fLDtC",
		},
	}

	for _, c := range tt {
		result := UrlHasher(c.user, c.url)
		assert.Equalf(t, c.short, result, "UrlHasher(%v, %s) should be %s, but is %s", c.user, c.url, c.short, result)
	}
}
