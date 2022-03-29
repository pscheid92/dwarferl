package internal

import "errors"

type UrlShortenerService struct {
	hasher    Hasher
	redirects RedirectRepository
}

func NewUrlShortenerService(hasher Hasher, redirects RedirectRepository) UrlShortenerService {
	return UrlShortenerService{
		hasher:    hasher,
		redirects: redirects,
	}
}

func (u UrlShortenerService) ShortenURL(url string) (string, error) {
	short := u.hasher(url)

	if err := u.redirects.Save(short, url); err != nil {
		return "", err
	}

	return short, nil
}

func (u UrlShortenerService) ExpandShortURL(short string) (string, error) {
	url, ok := u.redirects.Expand(short)
	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}

func (u UrlShortenerService) DeleteShortURL(short string) error {
	return u.redirects.Delete(short)
}
