-- name: GetUserByGoogleId :one
select * from users where google_provider_id = $1;

-- name: SaveUser :exec
INSERT INTO users (id, email, google_provider_id)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE SET email = excluded.email, google_provider_id = excluded.google_provider_id;
