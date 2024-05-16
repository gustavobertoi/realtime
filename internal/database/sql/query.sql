-- name: FindManyChannels :many
SELECT * FROM channels OFFSET $1 LIMIT $2;

-- name: FindOneChannelById :one
SELECT * FROM channels WHERE id = $1;

-- name: CreateChannel :exec
INSERT INTO channels (name, type, config, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateChannel :exec
UPDATE channels SET name = $1, type = $2, config = $3, updated_at = $4 WHERE id = $5;

-- name: DeleteChannel :exec
DELETE FROM channels WHERE id = $1;