package shortener

import (
	"context"
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

func (u UrlShortenerService) List(ctx context.Context, userID string) ([]internal.Redirect, error) {
	list, err := u.redirects.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (u UrlShortenerService) GetRedirectByShort(ctx context.Context, short string, userID string) (internal.Redirect, error) {
	redirect, err := u.redirects.GetRedirectByShort(ctx, short, userID)
	if err != nil {
		return internal.Redirect{}, err
	}
	return redirect, nil
}

func (u UrlShortenerService) ShortenURL(ctx context.Context, url string, userID string) (internal.Redirect, error) {
	redirect := internal.Redirect{
		UserID:    userID,
		Short:     u.hasher(userID, url),
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := u.redirects.Save(ctx, redirect); err != nil {
		return redirect, err
	}
	return redirect, nil
}

func (u UrlShortenerService) ExpandShortURL(ctx context.Context, short string) (string, error) {
	redirect, err := u.redirects.Expand(ctx, short)
	if err != nil {
		return "", err
	}
	return redirect, nil
}

func (u UrlShortenerService) DeleteShortURL(ctx context.Context, short string, userID string) error {
	return u.redirects.Delete(ctx, short, userID)
}
