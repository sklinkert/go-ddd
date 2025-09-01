-- name: CreateSeller :one
INSERT INTO sellers (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSellerById :one
SELECT id, name, created_at, updated_at
FROM sellers
WHERE id = $1;

-- name: GetAllSellers :many
SELECT id, name, created_at, updated_at
FROM sellers
ORDER BY created_at DESC;

-- name: UpdateSeller :exec
UPDATE sellers 
SET name = $2, updated_at = $3
WHERE id = $1;

-- name: DeleteSeller :exec
DELETE FROM sellers WHERE id = $1;