-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
$1,
$2,
$3,
$4,
$5
)
RETURNING *;

-- name: GetTokenByToken :one
SELECT
    *
FROM
    refresh_tokens rt
WHERE
    rt.Token = $1;


-- name: GetUserFromToken :one
SELECT
    rt.token,
    u.id,
    rt.expires_at,
    rt.revoked_AT
FROM
    refresh_tokens rt
JOIN
    users u on u.id = rt.user_id
WHERE
    rt.token = $1;


-- name: RevokeRefreshToken :one
UPDATE
    refresh_tokens
SET
    updated_at = $1,
    revoked_AT = $2
WHERE
    token = $3
RETURNING
*;
