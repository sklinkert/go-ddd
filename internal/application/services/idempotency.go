package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/domain/repositories"
)

// ErrRequestInFlight is returned when a request with the same idempotency
// key is still being processed by another caller.
var ErrRequestInFlight = errors.New("a request with this idempotency key is already in progress")

// ErrIdempotencyKeyReuse is returned when an idempotency key is reused with
// a different request payload (or a different operation).
var ErrIdempotencyKeyReuse = errors.New("idempotency key was already used with a different request")

// reservationTTL bounds how long a reservation without a stored response is
// honored. If the process crashes between reserving the key and storing the
// response, the next retry after the TTL takes the reservation over instead
// of returning 409 forever.
const reservationTTL = time.Minute

// withIdempotency wraps a command execution with idempotency handling:
//
//  1. The key is reserved atomically (INSERT .. ON CONFLICT DO NOTHING),
//     so concurrent requests with the same key cannot both execute.
//  2. A completed request with the same key and payload returns its cached
//     response; the same key with a different payload is rejected.
//  3. On failure the reservation is released so the client can retry;
//     reservations orphaned by a crash expire after reservationTTL.
func withIdempotency[T any](
	ctx context.Context,
	repo repositories.IdempotencyRepository,
	key string,
	cmd any,
	execute func() (*T, error),
) (*T, error) {
	if key == "" {
		return execute()
	}

	requestJSON, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("marshal idempotency request: %w", err)
	}

	record := entities.NewIdempotencyRecord(key, string(requestJSON))

	reserved := false
	for attempt := 0; attempt < 3 && !reserved; attempt++ {
		reserved, err = repo.Reserve(ctx, record)
		if err != nil {
			return nil, err
		}
		if reserved {
			break
		}

		existing, err := repo.FindByKey(ctx, key)
		if err != nil {
			return nil, err
		}
		if existing == nil {
			// Released between Reserve and FindByKey; try again.
			continue
		}

		if existing.Request != string(requestJSON) {
			return nil, ErrIdempotencyKeyReuse
		}

		if existing.IsCompleted() {
			var result T
			if err := json.Unmarshal([]byte(existing.Response), &result); err != nil {
				return nil, fmt.Errorf("unmarshal cached idempotency response: %w", err)
			}
			return &result, nil
		}

		if time.Since(existing.CreatedAt) < reservationTTL {
			return nil, ErrRequestInFlight
		}

		// Stale reservation: the previous holder crashed before completing.
		// Release it and retry the reservation.
		if err := repo.Delete(ctx, key); err != nil {
			return nil, err
		}
	}

	if !reserved {
		return nil, ErrRequestInFlight
	}

	result, err := execute()
	if err != nil {
		// Detached context: the failure may be a cancelled request context,
		// and the release must still reach the database.
		if deleteErr := repo.Delete(context.WithoutCancel(ctx), key); deleteErr != nil {
			slog.WarnContext(ctx, "failed to release idempotency key",
				slog.String("idempotency_key", key), slog.Any("error", deleteErr))
		}
		return nil, err
	}

	storeResponse(context.WithoutCancel(ctx), repo, key, result)

	return result, nil
}

// storeResponse persists the result against the reserved key. The write is
// best-effort: a failure is logged but must not fail the (already
// successful) operation.
func storeResponse(ctx context.Context, repo repositories.IdempotencyRepository, key string, result any) {
	responseJSON, err := json.Marshal(result)
	if err != nil {
		slog.WarnContext(ctx, "failed to marshal idempotency response",
			slog.String("idempotency_key", key), slog.Any("error", err))
		return
	}

	if err := repo.SetResponse(ctx, key, string(responseJSON), 200); err != nil {
		slog.WarnContext(ctx, "failed to persist idempotency response",
			slog.String("idempotency_key", key), slog.Any("error", err))
	}
}
