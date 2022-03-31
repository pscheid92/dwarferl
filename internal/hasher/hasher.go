package hasher

import (
	"github.com/jxskiss/base62"
	"github.com/pscheid92/dwarferl/internal"
	"hash/fnv"
)

func UrlHasher(user internal.User, url string) string {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(url))
	_, _ = hash.Write([]byte(user.ID.String()))
	hashed := hash.Sum([]byte{})

	return base62.EncodeToString(hashed)
}
