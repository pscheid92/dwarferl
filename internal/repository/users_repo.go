package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/repository/database"
)

type DBUsersRepository struct {
	queries *database.Queries
}

func NewDBUsersRepository(pool *pgxpool.Pool) *DBUsersRepository {
	return &DBUsersRepository{queries: database.New(pool)}
}

func (d *DBUsersRepository) Get(id string) (internal.User, error) {
	userDTO, err := d.queries.GetUserById(context.Background(), id)
	if err != nil {
		return internal.User{}, err
	}

	return internal.User{
		ID:    userDTO.ID,
		Email: userDTO.Email,
	}, nil
}
