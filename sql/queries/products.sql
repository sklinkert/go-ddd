-- name: CreateProduct :one
INSERT INTO products (id, name, price_cents, currency, seller_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProductById :one
SELECT p.id, p.name, p.price_cents, p.currency, p.seller_id, p.created_at, p.updated_at
FROM products p
JOIN sellers s ON p.seller_id = s.id
WHERE p.id = $1 AND p.deleted_at IS NULL AND s.deleted_at IS NULL;

-- name: GetAllProducts :many
SELECT p.id, p.name, p.price_cents, p.currency, p.seller_id, p.created_at, p.updated_at
FROM products p
JOIN sellers s ON p.seller_id = s.id
WHERE p.deleted_at IS NULL AND s.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: UpdateProduct :execrows
UPDATE products
SET name = $2, price_cents = $3, currency = $4, seller_id = $5, updated_at = $6
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteProduct :exec
UPDATE products SET deleted_at = NOW() WHERE id = $1;
