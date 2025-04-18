-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: RetrieveAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: RetrieveSelectChirp :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteSelectChirp :exec
DELETE FROM chirps WHERE id = $1;

-- name: RetrieveChirpsByAuthor :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at;
