package repository

import (
	"github.com/pscheid92/dwarferl/internal"
)

type StaticUsersRepository struct {
	// EMPTY
}

func (s StaticUsersRepository) Get(_ string) (internal.User, error) {
	user := internal.User{
		ID:    uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Email: "example@example.com",
	}
	return user, nil
}
