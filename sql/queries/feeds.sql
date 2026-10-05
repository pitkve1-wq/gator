-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name , url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
*
FROM
feeds;

-- name: GetUrl :one
SELECT
*
from
feeds
WHERE
$1 = url;

-- name: MarkFeedFetched :one
UPDATE feeds
SEt last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT
*
FROM
feeds
ORDER BY 
last_fetched_at NULLS FIRST
LIMIT 
1;

