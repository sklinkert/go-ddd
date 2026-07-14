-- Transactional outbox: domain events are stored in the same database as
-- the state change. A relay polls unpublished rows and forwards them to a
-- message broker, guaranteeing at-least-once delivery without dual writes.
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    event_name TEXT NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    published_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_outbox_events_unpublished ON outbox_events(occurred_at) WHERE published_at IS NULL;
