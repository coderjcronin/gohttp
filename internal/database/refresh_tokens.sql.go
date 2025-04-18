// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: refresh_tokens.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    NOW() + INTERVAL '60 days',
    $2
)
RETURNING token, created_at, updated_at, expires_at, revoked_at, user_id
`

type CreateRefreshTokenParams struct {
	Token  string
	UserID uuid.UUID
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const deleteAllTokens = `-- name: DeleteAllTokens :exec
DELETE FROM refresh_tokens
`

func (q *Queries) DeleteAllTokens(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllTokens)
	return err
}

const retrieveAllTokens = `-- name: RetrieveAllTokens :many
SELECT token, created_at, updated_at, expires_at, revoked_at, user_id FROM refresh_tokens ORDER BY created_at
`

func (q *Queries) RetrieveAllTokens(ctx context.Context) ([]RefreshToken, error) {
	rows, err := q.db.QueryContext(ctx, retrieveAllTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RefreshToken
	for rows.Next() {
		var i RefreshToken
		if err := rows.Scan(
			&i.Token,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExpiresAt,
			&i.RevokedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const retrieveSelectRefreshToken = `-- name: RetrieveSelectRefreshToken :one
SELECT token, created_at, updated_at, expires_at, revoked_at, user_id FROM refresh_tokens WHERE token = $1
`

func (q *Queries) RetrieveSelectRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, retrieveSelectRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET updated_at = NOW(), revoked_at = NOW() WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.ExecContext(ctx, revokeRefreshToken, token)
	return err
}
