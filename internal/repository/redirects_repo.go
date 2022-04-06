package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/pscheid92/dwarferl/internal/repository/database"
)

type DBRedirectsRepository struct {
	queries *database.Queries
}

func NewDBRedirectsRepository(pool *pgxpool.Pool) *DBRedirectsRepository {
	return &DBRedirectsRepository{queries: database.New(pool)}
}

func (d DBRedirectsRepository) List(userID string) ([]internal.Redirect, error) {
	redirectsDTO, err := d.queries.ListRedirectsByUserId(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	redirects := make([]internal.Redirect, len(redirectsDTO))
	for i, r := range redirectsDTO {
		redirects[i] = internal.Redirect{
			Short:     r.Short,
			URL:       r.Url,
			UserID:    r.UserID,
			CreatedAt: r.CreatedAt,
		}
	}

	return redirects, nil
}

func (d DBRedirectsRepository) Save(redirect internal.Redirect) error {
	params := database.SaveRedirectParams{
		Short:     redirect.Short,
		Url:       redirect.URL,
		UserID:    redirect.UserID,
		CreatedAt: redirect.CreatedAt,
	}
	err := d.queries.SaveRedirect(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (d DBRedirectsRepository) Expand(short string) (string, error) {
	url, err := d.queries.ExpandRedirect(context.Background(), short)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d DBRedirectsRepository) Delete(short string, userID string) error {
	return d.queries.DeleteRedirect(context.Background(), database.DeleteRedirectParams{
		Short:  short,
		UserID: userID,
	})
}
