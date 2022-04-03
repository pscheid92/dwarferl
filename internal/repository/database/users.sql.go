// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: users.sql

package database

import (
	"context"
)

const getUserById = `-- name: GetUserById :one
select id, email from users where id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(&i.ID, &i.Email)
	return i, err
}