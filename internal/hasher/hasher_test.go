package hasher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlHasher(t *testing.T) {
	tt := []struct {
		userID string
		url    string
		short  string
	}{
		{
			"00000000-0000-0000-0000-000000000000",
			"https://www.google.com",
			"hYhahA",
		},
		{
			"00000000-0000-0000-0000-000000000000",
			"https://patrickscheid.de",
			"5fLDtC",
		},
	}

	for _, c := range tt {
		result := UrlHasher(c.userID, c.url)
		assert.Equalf(t, c.short, result, "UrlHasher(%v, %s) should be %s, but is %s", c.userID, c.url, c.short, result)
	}
}
