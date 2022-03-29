package internal

import (
	"github.com/jxskiss/base62"
	"hash/fnv"
)

type Hasher = func(string) string

func UrlHasher(url string) string {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(url))
	hashed := hash.Sum32()

	bytes := base62.FormatUint(uint64(hashed))
	return string(bytes)
}
