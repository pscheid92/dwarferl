package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pscheid92/dwarferl/internal"
)

type DBRedirectsRepository struct {
	pool *pgxpool.Pool
}

func NewDBRedirectsRepository(pool *pgxpool.Pool) *DBRedirectsRepository {
	return &DBRedirectsRepository{pool: pool}
}

func (d DBRedirectsRepository) List(user internal.User) ([]internal.Redirect, error) {
	sql := `SELECT short, url, user_id, created_at FROM redirects WHERE user_id = $1`

	rows, err := d.pool.Query(context.Background(), sql, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var redirects []internal.Redirect
	for rows.Next() {
		var redirect internal.Redirect

		err = rows.Scan(&redirect.Short, &redirect.URL, &redirect.UserID, &redirect.CreatedAt)
		if err != nil {
			return nil, err
		}

		redirects = append(redirects, redirect)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return redirects, nil
}

func (d DBRedirectsRepository) Save(redirect internal.Redirect) error {
	// TODO(ps): we do not return the db entry, so creation time is plain wrong, if conflict!
	sql := `INSERT INTO redirects (short, url, user_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (short) DO NOTHING RE`

	_, err := d.pool.Exec(context.Background(), sql, redirect.Short, redirect.URL, redirect.UserID, redirect.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d DBRedirectsRepository) Expand(short string) (internal.Redirect, error) {
	sql := `SELECT short, url, user_id, created_at FROM redirects WHERE short = $1`

	var redirect internal.Redirect

	err := d.pool.
		QueryRow(context.Background(), sql, short).
		Scan(&redirect.Short, &redirect.URL, &redirect.UserID, &redirect.CreatedAt)

	if err != nil {
		return redirect, err
	}

	return redirect, nil
}

func (d DBRedirectsRepository) Delete(short string) error {
	sql := `DELETE FROM redirects WHERE short = $1`

	commandTag, err := d.pool.Exec(context.Background(), sql, short)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("failed to delete redirect")
	}

	return nil
}
