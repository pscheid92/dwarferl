package internal

import (
	"time"
)

type Hasher = func(User, string) string

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type Redirect struct {
	Short     string    `json:"short"`
	URL       string    `json:"url"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersRepository interface {
	Get(id string) (User, error)
}

type RedirectRepository interface {
	List(user User) ([]Redirect, error)
	Save(redirect Redirect) error
	Expand(short string) (string, error)
	Delete(short string) error
}

type UrlShortenerService interface {
	List(userID string) ([]Redirect, error)
	ShortenURL(url string) (Redirect, error)
	ExpandShortURL(short string) (string, error)
	DeleteShortURL(short string) error
}
