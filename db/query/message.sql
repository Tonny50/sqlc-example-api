-- name: CreateMessage :one
INSERT INTO message (thread_id, sender, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateThread :one
INSERT INTO thread (topic)
VALUES ($1)
RETURNING *;

-- name: GetMessageByID :one
SELECT * FROM message
WHERE id = $1;

-- name: GetMessagesByThread :many
SELECT * FROM message
WHERE thread_id = $1
ORDER BY created_at DESC;

-- name: UpdateMesageByID :one
UPDATE message
SET sender=$2, content=$3
WHERE id = $1
RETURNING *;

-- name: DeleteMeageByID :exec
DELETE FROM message
WHERE id = $1;