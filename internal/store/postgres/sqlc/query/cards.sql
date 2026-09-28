-- name: ListCards :many
SELECT c.*, creator.username AS creator_username, editor.username AS editor_username
FROM cards c
JOIN users creator ON creator.id = c.created_by
JOIN users editor ON editor.id = c.updated_by
ORDER BY c.created_at DESC, c.id;

-- name: FindCard :one
SELECT c.*, creator.username AS creator_username, editor.username AS editor_username
FROM cards c
JOIN users creator ON creator.id = c.created_by
JOIN users editor ON editor.id = c.updated_by
WHERE c.id = $1;

-- name: CreateCard :one
INSERT INTO cards (title, description, created_by, updated_by)
VALUES ($1, $2, sqlc.arg(user_id), sqlc.arg(user_id))
RETURNING id;

-- name: LockCard :one
SELECT version FROM cards WHERE id = $1 FOR UPDATE;

-- name: UpdateCard :execrows
UPDATE cards
SET title = $1, description = $2, status = $3, updated_by = $4,
    updated_at = now(), version = version + 1
WHERE id = $5 AND version = $6;

-- name: DeleteCard :execrows
DELETE FROM cards WHERE id = $1 AND version = $2;
