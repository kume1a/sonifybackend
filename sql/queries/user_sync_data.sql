-- name: CreateUserSyncDatum :one
INSERT INTO user_sync_data (
  id,
  user_id, 
  spotify_last_synced_at,
  user_audio_last_synced_at
) VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: GetUserSyncDatumByUserID :one
SELECT * FROM user_sync_data WHERE user_id = $1;

-- name: UpdateUserSyncDatumByUserID :one
UPDATE user_sync_data
  SET spotify_last_synced_at = COALESCE($1, spotify_last_synced_at),
      user_audio_last_synced_at = COALESCE($2, user_audio_last_synced_at)
WHERE user_id = $3
RETURNING *;
