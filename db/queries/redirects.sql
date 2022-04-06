-- name: ListRedirectsByUserId :many
SELECT *
FROM redirects
WHERE user_id = $1;

-- name: SaveRedirect :exec
INSERT INTO redirects (short, url, user_id, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (short) DO NOTHING
RETURNING *;

-- name: ExpandRedirect :one
SELECT url
FROM redirects
WHERE short = $1;

-- name: DeleteRedirect :exec
DELETE FROM redirects
WHERE short = $1 and user_id = $2;
