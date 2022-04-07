package internal

import (
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
	Save(user User) error
	GetByGoogleID(googleID string) (User, error)
}

type RedirectRepository interface {
	List(userID string) ([]Redirect, error)
	GetRedirectByShort(short string, userID string) (Redirect, error)
	Save(redirect Redirect) error
	Expand(short string) (string, error)
	Delete(short string, userID string) error
}

type UrlShortenerService interface {
	List(userID string) ([]Redirect, error)
	GetRedirectByShort(short string, userID string) (Redirect, error)
	ShortenURL(url string, userID string) (Redirect, error)
	ExpandShortURL(short string) (string, error)
	DeleteShortURL(short string, userID string) error
}

type UsersService interface {
	CreateWithGoogleID(googleID string, email string) (User, error)
	GetOrCreateByGoogle(googleID string, email string) (User, error)
}
