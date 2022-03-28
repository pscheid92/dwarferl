package internal

import (
	"errors"
	"github.com/jxskiss/base62"
	"hash/fnv"
)

type UrlShortener struct {
	storage Storage
}

func NewUrlShortener(storage Storage) UrlShortener {
	return UrlShortener{storage: storage}
}

func (u UrlShortener) Shorten(url string) (string, error) {
	hash := fnv.New32a()
	hash.Write([]byte(url))
	hashed := hash.Sum32()

	bytes := base62.FormatUint(uint64(hashed))
	short := string(bytes)

	err := u.storage.SaveRedirect(short, url)
	if err != nil {
		return "", err
	}

	return short, nil
}

func (u UrlShortener) Expand(short string) (string, error) {
	url, ok := u.storage.ExpandRedirect(short)
	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}

func (u UrlShortener) Delete(short string) error {
	return u.storage.DeleteRedirect(short)
}
