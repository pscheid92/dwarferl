-- name: GetUserById :one
select * from users where id = $1;
