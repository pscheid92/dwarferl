package internal

import (
	"context"
	"time"
)

type User struct {
	ID       string
	Email    string
	GoogleID string
}

type Redirect struct {
	Short     string
	URL       string
	UserID    string
	CreatedAt time.Time
}

type Hasher interface {
	Hash(userID string, url string) string
	Validate(short string) bool
}

type UsersRepository interface {
	Save(ctx context.Context, user User) error
	GetByGoogleID(ctx context.Context, googleID string) (User, error)
}

type RedirectRepository interface {
	List(ctx context.Context, userID string) ([]Redirect, error)
	GetRedirectByShort(ctx context.Context, short string, userID string) (Redirect, error)
	Save(ctx context.Context, redirect Redirect) error
	Expand(ctx context.Context, short string) (string, error)
	Delete(ctx context.Context, short string, userID string) error
}

type UrlShortenerService interface {
	List(ctx context.Context, userID string) ([]Redirect, error)
	GetRedirectByShort(ctx context.Context, short string, userID string) (Redirect, error)
	ShortenURL(ctx context.Context, url string, userID string) (Redirect, error)
	ExpandShortURL(ctx context.Context, short string) (string, error)
	DeleteShortURL(ctx context.Context, short string, userID string) error
}

type UsersService interface {
	CreateWithGoogleID(ctx context.Context, googleID string, email string) (User, error)
	GetOrCreateByGoogle(ctx context.Context, googleID string, email string) (User, error)
}
