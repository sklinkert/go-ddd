-- name: CreateIdempotencyRecord :one
INSERT INTO idempotency_records (id, key, request, response, status_code, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetIdempotencyRecordByKey :one
SELECT id, key, request, response, status_code, created_at
FROM idempotency_records
WHERE key = $1;

-- name: UpdateIdempotencyRecord :one
UPDATE idempotency_records
SET request = $2, response = $3, status_code = $4
WHERE id = $1
RETURNING *;