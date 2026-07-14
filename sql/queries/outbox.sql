-- name: InsertOutboxEvent :exec
INSERT INTO outbox_events (id, aggregate_id, event_name, payload, occurred_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetUnpublishedOutboxEvents :many
SELECT id, aggregate_id, event_name, payload, occurred_at, published_at
FROM outbox_events
WHERE published_at IS NULL
ORDER BY occurred_at
LIMIT $1;

-- name: MarkOutboxEventPublished :exec
UPDATE outbox_events SET published_at = NOW() WHERE id = $1;
