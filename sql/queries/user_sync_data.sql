-- name: CreateUserSyncData :one
INSERT INTO user_sync_data (
  id,
  user_id, 
  spotify_last_synced_at
) VALUES ($1,$2,$3)
RETURNING *;

-- name: GetUserSyncDatumByUserId :one
SELECT * FROM user_sync_data WHERE user_id = $1;

-- name: UpdateUserSyncDatumByUserId :one
UPDATE user_sync_data
  SET spotify_last_synced_at = COALESCE($1, spotify_last_synced_at)
  WHERE user_id = $2
  RETURNING *;
