package postgres

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

// The minor-units migration rewrites pending product.created payloads from
// PriceCents to PriceMinorUnits. Without the rewrite, a consumer of the new
// schema would read a silent zero price from pre-rename outbox rows. The
// test seeds a legacy payload and re-runs the migration's UPDATE from the
// actual migration file, so the two cannot drift apart.
func TestMigration000004_RewritesPendingOutboxPayloads(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	ctx := context.Background()

	legacyPayload := `{"Id":"0198c0de-0000-7000-8000-000000000001","Aggregate":"0198c0de-0000-7000-8000-000000000002","OccurredAtT":"2026-07-01T00:00:00Z","Name":"Legacy Product","PriceCents":4999,"Currency":"EUR","SellerId":"0198c0de-0000-7000-8000-000000000003"}`
	eventId := uuid.Must(uuid.NewV7())
	require.NoError(t, testDB.Queries.InsertOutboxEvent(ctx, db.InsertOutboxEventParams{
		ID:          eventId,
		AggregateID: uuid.Must(uuid.NewV7()),
		EventName:   "product.created",
		Payload:     []byte(legacyPayload),
		OccurredAt:  timestamptzFromTime(time.Now()),
	}))

	migration, err := os.ReadFile(filepath.Join(testhelpers.ProjectRoot(t), "migrations", "000004_price_minor_units.up.sql"))
	require.NoError(t, err)
	// The column rename already ran during schema setup; re-running just the
	// payload UPDATE is idempotent thanks to its `payload ? 'PriceCents'` guard.
	_, err = testDB.Pool.Exec(ctx, extractOutboxUpdate(t, string(migration)))
	require.NoError(t, err)

	events, err := testDB.Queries.GetUnpublishedOutboxEvents(ctx, 10)
	require.NoError(t, err)
	require.Len(t, events, 1)

	var payload map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(events[0].Payload, &payload))

	assert.NotContains(t, payload, "PriceCents")
	require.Contains(t, payload, "PriceMinorUnits")
	assert.Equal(t, "4999", string(payload["PriceMinorUnits"]))
	assert.Equal(t, `"Legacy Product"`, string(payload["Name"]))
}

// extractOutboxUpdate pulls the UPDATE statement out of the migration file
// so the test always runs what the migration actually ships.
func extractOutboxUpdate(t *testing.T, migration string) string {
	t.Helper()

	idx := strings.Index(migration, "UPDATE")
	require.GreaterOrEqual(t, idx, 0, "migration must contain the outbox UPDATE statement")

	return migration[idx:]
}
