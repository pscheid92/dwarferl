package internal

import "github.com/google/uuid"

type Hasher = func(User, string) string

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UsersRepository interface {
	Get() (User, error)
}

type RedirectRepository interface {
	Save(short string, url string) error
	Expand(short string) (string, bool)
	Delete(short string) error
}

type UrlShortenerService interface {
	ShortenURL(url string) (string, error)
	ExpandShortURL(short string) (string, error)
	DeleteShortURL(short string) error
}
