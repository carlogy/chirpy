-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
$1,
$2,
$3,
$4,
$5
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT
    *
FROM
    users u
WHERE
    u.email = $1;


-- name: UpdateUserDetails :one
UPDATE
    users
SET
    updated_at = $1,
    email = $2,
    hashed_password = $3
WHERE
    id = $4
RETURNING *;

-- name: UpgradeUserToRed :one
UPDATE
    users
SET
    updated_at = $1,
    is_chirpy_red = $2
WHERE
    id = $3
RETURNING *;
