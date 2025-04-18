-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: QueryHashedPassword :one
SELECT * FROM users WHERE email = $1;

-- name: RetrieveBasedOnID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserPassword :one
UPDATE users SET hashed_password = $1, email = $2, updated_at = NOW() WHERE id = $3 RETURNING users.*; 

-- name: UpgradeToChirpyRed :one
UPDATE users SET is_chirpy_red = true, updated_at = NOW() WHERE id = $1 RETURNING users.*;