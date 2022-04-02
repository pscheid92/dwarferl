package hasher

import (
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlHasher(t *testing.T) {
	var testUser = internal.User{
		ID:    "00000000-0000-0000-0000-000000000000",
		Email: "example@example.com",
	}

	tt := []struct {
		user  internal.User
		url   string
		short string
	}{
		{
			testUser,
			"https://www.google.com",
			"hYhahA",
		},
		{
			testUser,
			"https://patrickscheid.de",
			"5fLDtC",
		},
	}

	for _, c := range tt {
		result := UrlHasher(c.user, c.url)
		assert.Equalf(t, c.short, result, "UrlHasher(%v, %s) should be %s, but is %s", c.user, c.url, c.short, result)
	}
}
