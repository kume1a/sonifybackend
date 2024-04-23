// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: audio_like.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createAudioLike = `-- name: CreateAudioLike :one
INSERT INTO audio_likes (
  audio_id, 
  user_id
) VALUES ($1, $2) 
RETURNING user_id, audio_id
`

type CreateAudioLikeParams struct {
	AudioID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) CreateAudioLike(ctx context.Context, arg CreateAudioLikeParams) (AudioLike, error) {
	row := q.db.QueryRowContext(ctx, createAudioLike, arg.AudioID, arg.UserID)
	var i AudioLike
	err := row.Scan(&i.UserID, &i.AudioID)
	return i, err
}

const deleteAudioLike = `-- name: DeleteAudioLike :exec
DELETE FROM audio_likes 
  WHERE audio_id = $1 AND user_id = $2
`

type DeleteAudioLikeParams struct {
	AudioID uuid.UUID
	UserID  uuid.UUID
}

func (q *Queries) DeleteAudioLike(ctx context.Context, arg DeleteAudioLikeParams) error {
	_, err := q.db.ExecContext(ctx, deleteAudioLike, arg.AudioID, arg.UserID)
	return err
}

const deleteUserAudioLikesByAudioIDs = `-- name: DeleteUserAudioLikesByAudioIDs :exec
DELETE FROM audio_likes WHERE user_id = $1 AND audio_id = ANY($2::uuid[])
`

type DeleteUserAudioLikesByAudioIDsParams struct {
	UserID   uuid.UUID
	AudioIds []uuid.UUID
}

func (q *Queries) DeleteUserAudioLikesByAudioIDs(ctx context.Context, arg DeleteUserAudioLikesByAudioIDsParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserAudioLikesByAudioIDs, arg.UserID, pq.Array(arg.AudioIds))
	return err
}

const getAudioLikesByUserId = `-- name: GetAudioLikesByUserId :many
SELECT user_id, audio_id FROM audio_likes WHERE user_id = $1
`

func (q *Queries) GetAudioLikesByUserId(ctx context.Context, userID uuid.UUID) ([]AudioLike, error) {
	rows, err := q.db.QueryContext(ctx, getAudioLikesByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AudioLike
	for rows.Next() {
		var i AudioLike
		if err := rows.Scan(&i.UserID, &i.AudioID); err != nil {
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
