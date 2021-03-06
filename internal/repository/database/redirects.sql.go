// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: redirects.sql

package database

import (
	"context"
	"time"
)

const deleteRedirect = `-- name: DeleteRedirect :exec
DELETE FROM redirects
WHERE short = $1 and user_id = $2
`

type DeleteRedirectParams struct {
	Short  string
	UserID string
}

func (q *Queries) DeleteRedirect(ctx context.Context, arg DeleteRedirectParams) error {
	_, err := q.db.Exec(ctx, deleteRedirect, arg.Short, arg.UserID)
	return err
}

const expandRedirect = `-- name: ExpandRedirect :one
SELECT url
FROM redirects
WHERE short = $1
`

func (q *Queries) ExpandRedirect(ctx context.Context, short string) (string, error) {
	row := q.db.QueryRow(ctx, expandRedirect, short)
	var url string
	err := row.Scan(&url)
	return url, err
}

const getRedirectByShort = `-- name: GetRedirectByShort :one
SELECT short, url, user_id, created_at
FROM redirects
WHERE short = $1 and user_id = $2
`

type GetRedirectByShortParams struct {
	Short  string
	UserID string
}

func (q *Queries) GetRedirectByShort(ctx context.Context, arg GetRedirectByShortParams) (Redirect, error) {
	row := q.db.QueryRow(ctx, getRedirectByShort, arg.Short, arg.UserID)
	var i Redirect
	err := row.Scan(
		&i.Short,
		&i.Url,
		&i.UserID,
		&i.CreatedAt,
	)
	return i, err
}

const listRedirectsByUserId = `-- name: ListRedirectsByUserId :many
SELECT short, url, user_id, created_at
FROM redirects
WHERE user_id = $1
`

func (q *Queries) ListRedirectsByUserId(ctx context.Context, userID string) ([]Redirect, error) {
	rows, err := q.db.Query(ctx, listRedirectsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Redirect
	for rows.Next() {
		var i Redirect
		if err := rows.Scan(
			&i.Short,
			&i.Url,
			&i.UserID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveRedirect = `-- name: SaveRedirect :exec
INSERT INTO redirects (short, url, user_id, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short) DO NOTHING
RETURNING short, url, user_id, created_at
`

type SaveRedirectParams struct {
	Short     string
	Url       string
	UserID    string
	CreatedAt time.Time
}

func (q *Queries) SaveRedirect(ctx context.Context, arg SaveRedirectParams) error {
	_, err := q.db.Exec(ctx, saveRedirect,
		arg.Short,
		arg.Url,
		arg.UserID,
		arg.CreatedAt,
	)
	return err
}
