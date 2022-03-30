package internal

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type UsersRepository interface {
	Get() (User, error)
}

type StaticUsersRepository struct {
	// EMPTY
}

var DummyUser = User{
	ID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	Email: "example@example.com",
}

func (s StaticUsersRepository) Get() (User, error) {
	return DummyUser, nil
}
