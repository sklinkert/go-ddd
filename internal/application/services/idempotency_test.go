package services

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResult struct {
	Value string `json:"value"`
}

func TestWithIdempotency_EmptyKeyBypasses(t *testing.T) {
	repo := NewMockIdempotencyRepository()

	executions := 0
	result, err := withIdempotency(context.Background(), repo, "", "cmd", func() (*testResult, error) {
		executions++
		return &testResult{Value: "fresh"}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "fresh", result.Value)
	assert.Equal(t, 1, executions)
	assert.Zero(t, repo.reserveCalls)
	assert.Empty(t, repo.records)
}

func TestWithIdempotency_CachedCompletedResponseReturned(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	record := entities.NewIdempotencyRecord("key-1", `"cmd"`)
	record.SetResponse(`{"value":"cached"}`, 200)
	repo.records["key-1"] = record

	executions := 0
	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		executions++
		return &testResult{Value: "fresh"}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "cached", result.Value)
	assert.Zero(t, executions, "cached response must not re-execute the command")
}

func TestWithIdempotency_InFlightRecordReturnsErrRequestInFlight(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	repo.records["key-1"] = entities.NewIdempotencyRecord("key-1", `"cmd"`) // reserved, no response yet

	executions := 0
	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		executions++
		return &testResult{Value: "fresh"}, nil
	})

	assert.ErrorIs(t, err, ErrRequestInFlight)
	assert.Nil(t, result)
	assert.Zero(t, executions)
}

func TestWithIdempotency_ReservationReleasedOnExecuteError(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	executeErr := errors.New("boom")

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return nil, executeErr
	})

	assert.ErrorIs(t, err, executeErr)
	assert.Nil(t, result)
	assert.Equal(t, 1, repo.deleteCalls)
	assert.Equal(t, []string{"key-1"}, repo.deletedKeys)
	assert.Empty(t, repo.records, "failed execution must release the key so the client can retry")
}

func TestWithIdempotency_SuccessStoresResponse(t *testing.T) {
	repo := NewMockIdempotencyRepository()

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return &testResult{Value: "fresh"}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "fresh", result.Value)

	record := repo.records["key-1"]
	require.NotNil(t, record)
	assert.True(t, record.IsCompleted())
	assert.Equal(t, 200, record.StatusCode)
	assert.JSONEq(t, `{"value":"fresh"}`, record.Response)
}

func TestWithIdempotency_ReserveErrorPropagates(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	repo.reserveErr = errors.New("db down")

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return &testResult{Value: "fresh"}, nil
	})

	assert.EqualError(t, err, "db down")
	assert.Nil(t, result)
}

func TestWithIdempotency_FindErrorPropagates(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	repo.records["key-1"] = entities.NewIdempotencyRecord("key-1", `"cmd"`) // key already reserved
	repo.findErr = errors.New("db down")

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return &testResult{Value: "fresh"}, nil
	})

	assert.EqualError(t, err, "db down")
	assert.Nil(t, result)
}

// raceLosingRepo simulates losing the reservation race: FindByKey sees no
// record, but Reserve fails because a concurrent request claimed the key.
type raceLosingRepo struct {
	*MockIdempotencyRepository
	firstFind bool
}

func (r *raceLosingRepo) FindByKey(ctx context.Context, key string) (*entities.IdempotencyRecord, error) {
	if !r.firstFind {
		r.firstFind = true
		return nil, nil
	}
	return r.MockIdempotencyRepository.FindByKey(ctx, key)
}

func (r *raceLosingRepo) Reserve(ctx context.Context, record *entities.IdempotencyRecord) (bool, error) {
	return false, nil
}

func TestWithIdempotency_LostRaceServesWinnersResponse(t *testing.T) {
	inner := NewMockIdempotencyRepository()
	record := entities.NewIdempotencyRecord("key-1", `"cmd"`)
	record.SetResponse(`{"value":"winner"}`, 200)
	inner.records["key-1"] = record

	repo := &raceLosingRepo{MockIdempotencyRepository: inner}

	executions := 0
	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		executions++
		return &testResult{Value: "loser"}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "winner", result.Value)
	assert.Zero(t, executions)
}

func TestWithIdempotency_LostRaceStillInFlight(t *testing.T) {
	inner := NewMockIdempotencyRepository()
	inner.records["key-1"] = entities.NewIdempotencyRecord("key-1", `"cmd"`) // winner not finished yet

	repo := &raceLosingRepo{MockIdempotencyRepository: inner}

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return &testResult{Value: "loser"}, nil
	})

	assert.ErrorIs(t, err, ErrRequestInFlight)
	assert.Nil(t, result)
}

func TestWithIdempotency_ConcurrentSameKeyExecutesOnce(t *testing.T) {
	repo := NewMockIdempotencyRepository()

	var mu sync.Mutex
	executions := 0
	successes := 0
	inFlight := 0

	const callers = 8
	var wg sync.WaitGroup
	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := withIdempotency(context.Background(), repo, "shared-key", "cmd", func() (*testResult, error) {
				mu.Lock()
				executions++
				mu.Unlock()
				return &testResult{Value: "fresh"}, nil
			})
			mu.Lock()
			defer mu.Unlock()
			switch {
			case err == nil:
				successes++
			case errors.Is(err, ErrRequestInFlight):
				inFlight++
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, executions, "exactly one caller may execute")
	assert.Equal(t, callers, successes+inFlight)
	assert.GreaterOrEqual(t, successes, 1)
}

func TestWithIdempotency_KeyReuseWithDifferentPayloadRejected(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	record := entities.NewIdempotencyRecord("key-1", `"other-cmd"`)
	record.SetResponse(`{"value":"cached"}`, 200)
	repo.records["key-1"] = record

	executions := 0
	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		executions++
		return &testResult{Value: "fresh"}, nil
	})

	assert.ErrorIs(t, err, ErrIdempotencyKeyReuse)
	assert.Nil(t, result)
	assert.Zero(t, executions)
}

func TestWithIdempotency_StaleReservationTakenOver(t *testing.T) {
	repo := NewMockIdempotencyRepository()
	stale := entities.NewIdempotencyRecord("key-1", `"cmd"`)
	stale.CreatedAt = stale.CreatedAt.Add(-2 * reservationTTL) // crashed holder
	repo.records["key-1"] = stale

	result, err := withIdempotency(context.Background(), repo, "key-1", "cmd", func() (*testResult, error) {
		return &testResult{Value: "fresh"}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "fresh", result.Value)
	record := repo.records["key-1"]
	require.NotNil(t, record)
	assert.True(t, record.IsCompleted())
}
