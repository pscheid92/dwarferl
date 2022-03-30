package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlHasher(t *testing.T) {
	tt := []struct {
		url   string
		short string
	}{
		{
			"https://www.google.com",
			"CdhQjT",
		},
		{
			"https://patrickscheid.de",
			"D7QWOh",
		},
	}

	for _, c := range tt {
		result := UrlHasher(c.url)
		assert.Equalf(t, c.short, result, "UrlHasher(%s) should be %s, but is %s", c.url, c.short, result)
	}
}
