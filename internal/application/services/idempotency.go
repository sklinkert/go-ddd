package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

// newIdempotencyRecord builds an idempotency record for the given command.
// It returns nil when no idempotency key was supplied, or when the command
// cannot be serialized (in which case idempotency is skipped rather than
// failing the operation).
func newIdempotencyRecord(ctx context.Context, key string, cmd any) *entities.IdempotencyRecord {
	if key == "" {
		return nil
	}

	requestJSON, err := json.Marshal(cmd)
	if err != nil {
		slog.WarnContext(ctx, "failed to marshal idempotency request; skipping idempotency",
			slog.String("idempotency_key", key), slog.Any("error", err))
		return nil
	}

	return entities.NewIdempotencyRecord(key, string(requestJSON))
}

// storeIdempotencyResponse persists the operation result against the
// idempotency record. The write is best-effort: a failure is logged but must
// not fail the (already successful) operation.
func storeIdempotencyResponse(ctx context.Context, repo repositories.IdempotencyRepository, record *entities.IdempotencyRecord, result any) {
	if record == nil {
		return
	}

	responseJSON, err := json.Marshal(result)
	if err != nil {
		slog.WarnContext(ctx, "failed to marshal idempotency response; skipping idempotency write",
			slog.String("idempotency_key", record.Key), slog.Any("error", err))
		return
	}

	record.SetResponse(string(responseJSON), 200)
	if _, err := repo.Create(ctx, record); err != nil {
		slog.WarnContext(ctx, "failed to persist idempotency record",
			slog.String("idempotency_key", record.Key), slog.Any("error", err))
	}
}
