package shortener

import (
	"errors"
	"github.com/pscheid92/dwarferl/internal"
)

type UrlShortenerService struct {
	hasher    internal.Hasher
	redirects internal.RedirectRepository
	users     internal.UsersRepository
}

func NewUrlShortenerService(hasher internal.Hasher, redirects internal.RedirectRepository, users internal.UsersRepository) UrlShortenerService {
	return UrlShortenerService{
		hasher:    hasher,
		redirects: redirects,
		users:     users,
	}
}

func (u UrlShortenerService) List(userID string) (map[string]string, error) {
	user, err := u.users.Get(userID)
	if err != nil {
		return nil, err
	}

	list, err := u.redirects.List(user)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u UrlShortenerService) ShortenURL(url string) (string, error) {
	user, _ := u.users.Get("")
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
