package internal

import (
	"github.com/jxskiss/base62"
	"hash/fnv"
)

type Hasher = func(User, string) string

func UrlHasher(user User, url string) string {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(url))
	_, _ = hash.Write([]byte(user.ID.String()))
	hashed := hash.Sum([]byte{})

	return base62.EncodeToString(hashed)
}
