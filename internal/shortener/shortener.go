package shortener

import (
	"github.com/pscheid92/dwarferl/internal"
	"time"
)

type UrlShortenerService struct {
	hasher    internal.Hasher
	redirects internal.RedirectRepository
}

func NewUrlShortenerService(hasher internal.Hasher, redirects internal.RedirectRepository) UrlShortenerService {
	return UrlShortenerService{
		hasher:    hasher,
		redirects: redirects,
	}
}

func (u UrlShortenerService) List(userID string) ([]internal.Redirect, error) {
	list, err := u.redirects.List(userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u UrlShortenerService) GetRedirectByShort(short string, userID string) (internal.Redirect, error) {
	redirect, err := u.redirects.GetRedirectByShort(short, userID)
	if err != nil {
		return internal.Redirect{}, err
	}
	return redirect, nil
}

func (u UrlShortenerService) ShortenURL(url string, userID string) (internal.Redirect, error) {
	redirect := internal.Redirect{
		UserID:    userID,
		Short:     u.hasher(userID, url),
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

func (u UrlShortenerService) DeleteShortURL(short string, userID string) error {
	return u.redirects.Delete(short, userID)
}
