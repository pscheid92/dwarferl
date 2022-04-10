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

func (d DBRedirectsRepository) List(ctx context.Context, userID string) ([]internal.Redirect, error) {
	dtos, err := d.queries.ListRedirectsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	redirects := make([]internal.Redirect, len(dtos))
	for i, r := range dtos {
		redirects[i] = dtoToRedirect(r)
	}

	return redirects, nil
}

func (d DBRedirectsRepository) GetRedirectByShort(ctx context.Context, short string, userID string) (internal.Redirect, error) {
	args := database.GetRedirectByShortParams{Short: short, UserID: userID}
	dto, err := d.queries.GetRedirectByShort(ctx, args)
	if err != nil {
		return internal.Redirect{}, err
	}
	return dtoToRedirect(dto), nil
}

func (d DBRedirectsRepository) Save(ctx context.Context, redirect internal.Redirect) error {
	params := database.SaveRedirectParams{
		Short:     redirect.Short,
		Url:       redirect.URL,
		UserID:    redirect.UserID,
		CreatedAt: redirect.CreatedAt,
	}
	err := d.queries.SaveRedirect(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (d DBRedirectsRepository) Expand(ctx context.Context, short string) (string, error) {
	url, err := d.queries.ExpandRedirect(ctx, short)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (d DBRedirectsRepository) Delete(ctx context.Context, short string, userID string) error {
	return d.queries.DeleteRedirect(ctx, database.DeleteRedirectParams{
		Short:  short,
		UserID: userID,
	})
}

func dtoToRedirect(dto database.Redirect) internal.Redirect {
	return internal.Redirect{
		Short:     dto.Short,
		URL:       dto.Url,
		UserID:    dto.UserID,
		CreatedAt: dto.CreatedAt,
	}
}
