-- name: CreateUser :one
INSERT INTO users (id , created_at , updated_at , name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT
*
from
users
WHERE
name = $1;

-- name: ResetUser :exec
DELETE from
users;

-- name: GetUsers :many
SELECT
*
from
users;

-- name: Getwhat :one
SELECT
name
from
users
WHERE
id = $1;


