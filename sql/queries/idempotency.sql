-- name: ReserveIdempotencyKey :execrows
-- Atomically claims the key. Zero rows means another request already holds it.
INSERT INTO idempotency_records (id, key, request, response, status_code, created_at)
VALUES ($1, $2, $3, '', 0, $4)
ON CONFLICT (key) DO NOTHING;

-- name: GetIdempotencyRecordByKey :one
SELECT id, key, request, response, status_code, created_at
FROM idempotency_records
WHERE key = $1;

-- name: SetIdempotencyResponse :exec
UPDATE idempotency_records
SET response = $2, status_code = $3
WHERE key = $1;

-- name: DeleteIdempotencyRecord :exec
DELETE FROM idempotency_records WHERE key = $1;
