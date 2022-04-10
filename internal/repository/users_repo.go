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

func (d *DBUsersRepository) Save(ctx context.Context, user internal.User) error {
	return d.queries.SaveUser(ctx, database.SaveUserParams{
		ID:               user.ID,
		Email:            user.Email,
		GoogleProviderID: user.GoogleID,
	})
}

func (d *DBUsersRepository) GetByGoogleID(ctx context.Context, googleID string) (internal.User, error) {
	userDTO, err := d.queries.GetUserByGoogleId(ctx, googleID)
	if err != nil {
		return internal.User{}, err
	}

	return dtoToUser(userDTO), nil
}

func dtoToUser(dto database.User) internal.User {
	return internal.User{
		ID:       dto.ID,
		Email:    dto.Email,
		GoogleID: dto.GoogleProviderID,
	}
}
