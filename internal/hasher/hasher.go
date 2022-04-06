package hasher

import (
	"github.com/jxskiss/base62"
	"hash/fnv"
)

func UrlHasher(userID string, url string) string {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(url))
	_, _ = hash.Write([]byte(userID))
	hashed := hash.Sum([]byte{})

	return base62.EncodeToString(hashed)
}
