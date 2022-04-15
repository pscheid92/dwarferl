package hasher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUrlHasher_Hash(t *testing.T) {
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

	hasher := NewUrlHasher()

	for _, c := range tt {
		result := hasher.Hash(c.userID, c.url)
		assert.Equalf(t, c.short, result, "hash of (%v, %s) should be %s, but is %s", c.userID, c.url, c.short, result)
	}
}

func TestUrlHasher_Validate(t *testing.T) {
	tt := []struct {
		short    string
		expected bool
	}{
		{"hYhahA", true},
		{"5fLDtC", true},
		{"", false},
		{"falséy", false},
		{"short", false},
		{"toolong", false},
	}

	hasher := NewUrlHasher()

	for _, c := range tt {
		result := hasher.Validate(c.short)
		assert.Equalf(t, c.expected, result, "validation of '%v' should be %t, but is %t", c.short, c.expected, result)
	}
}
