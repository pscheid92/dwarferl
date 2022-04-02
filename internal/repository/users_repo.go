package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pscheid92/dwarferl/internal"
)

type DBUsersRepository struct {
	pool *pgxpool.Pool
}

func NewDBUsersRepository(pool *pgxpool.Pool) *DBUsersRepository {
	return &DBUsersRepository{pool: pool}
}

func (d *DBUsersRepository) Get(id string) (internal.User, error) {
	sql := `SELECT id, email FROM users WHERE id = $1`

	var user internal.User
	if err := d.pool.QueryRow(context.Background(), sql, id).Scan(&user.ID, &user.Email); err != nil {
		return user, err
	}

	return user, nil
}
