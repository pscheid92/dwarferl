package internal

import "errors"

type UrlShortenerService struct {
	hasher    Hasher
	redirects RedirectRepository
	users     UsersRepository
}

func NewUrlShortenerService(hasher Hasher, redirects RedirectRepository, users UsersRepository) UrlShortenerService {
	return UrlShortenerService{
		hasher:    hasher,
		redirects: redirects,
		users:     users,
	}
}

func (u UrlShortenerService) ShortenURL(url string) (string, error) {
	user, _ := u.users.Get()
	short := u.hasher(user, url)

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
