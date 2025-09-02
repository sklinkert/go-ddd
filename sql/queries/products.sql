-- name: CreateProduct :one
INSERT INTO products (id, name, price, seller_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductById :one
SELECT p.id, p.name, p.price, p.seller_id, p.created_at, p.updated_at,
       s.id as s_id, s.name as s_name, s.created_at as s_created_at, s.updated_at as s_updated_at
FROM products p
JOIN sellers s ON p.seller_id = s.id
WHERE p.id = $1 AND p.deleted_at IS NULL AND s.deleted_at IS NULL;

-- name: GetAllProducts :many
SELECT p.id, p.name, p.price, p.seller_id, p.created_at, p.updated_at,
       s.id as s_id, s.name as s_name, s.created_at as s_created_at, s.updated_at as s_updated_at
FROM products p
JOIN sellers s ON p.seller_id = s.id
WHERE p.deleted_at IS NULL AND s.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: UpdateProduct :exec
UPDATE products 
SET name = $2, price = $3, seller_id = $4, updated_at = $5
WHERE id = $1;

-- name: DeleteProduct :exec
UPDATE products SET deleted_at = NOW() WHERE id = $1;