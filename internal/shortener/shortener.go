package shortener

import (
	"github.com/pscheid92/dwarferl/internal"
	"time"
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

func (u UrlShortenerService) List(userID string) ([]internal.Redirect, error) {
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

func (u UrlShortenerService) ShortenURL(url string) (internal.Redirect, error) {
	user, _ := u.users.Get("00000000-0000-0000-0000-000000000000")
	short := u.hasher(user, url)

	redirect := internal.Redirect{
		UserID:    user.ID,
		Short:     short,
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := u.redirects.Save(redirect); err != nil {
		return redirect, err
	}

	return redirect, nil
}

func (u UrlShortenerService) ExpandShortURL(short string) (string, error) {
	redirect, err := u.redirects.Expand(short)
	if err != nil {
		return "", err
	}
	return redirect, nil
}

func (u UrlShortenerService) DeleteShortURL(short string) error {
	return u.redirects.Delete(short)
}
