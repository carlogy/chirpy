-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
$1,
$2,
$3,
$4,
$5
)
RETURNING *;


-- name: GetChirps :many
SELECT
    *
FROM
    chirps c
ORDER BY
    c.created_at;

-- name: GetChirpsDesc :many
SELECT
    *
FROM
    chirps c
ORDER BY
    c.created_at DESC;


-- name: GetChirp :one
SELECT
    *
FROM
    chirps c
WHERE
    c.id = $1;

-- name: DeleteChirp :one
DELETE FROM
    chirps
WHERE
    id = $1
AND
    user_id = $2
RETURNING *;

-- name: GetChirpsByAuthor :many
SELECT
    *
FROM
    chirps c
WHERE
    c.user_id = $1
ORDER BY
    c.created_at;

-- name: GetChirpsByAuthorDESC :many
SELECT
    *
FROM
    chirps c
WHERE
    c.user_id = $1
ORDER BY
    c.created_at DESC;
