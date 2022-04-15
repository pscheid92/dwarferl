package hasher

import (
	"github.com/jxskiss/base62"
	"hash/fnv"
	"regexp"
)

type UrlHasher struct {
	regex *regexp.Regexp
}

func NewUrlHasher() UrlHasher {
	return UrlHasher{regex: regexp.MustCompile(`^[A-Za-z\d]{6}$`)}
}

func (UrlHasher) Hash(userID string, url string) string {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(url))
	_, _ = hash.Write([]byte(userID))
	hashed := hash.Sum([]byte{})
	return base62.EncodeToString(hashed)
}

func (h UrlHasher) Validate(short string) bool {
	return h.regex.MatchString(short)
}
