package internal

import "testing"

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
		if result != c.short {
			t.Errorf("UrlHasher(%s) == %s, want %s", c.url, result, c.short)
		}
	}
}
