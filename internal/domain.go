package internal

import (
	"context"
	"time"
)

type Hasher = func(string, string) string

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	GoogleID string `json:"google_id"`
}

type Redirect struct {
	Short     string    `json:"short"`
	URL       string    `json:"url"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
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
