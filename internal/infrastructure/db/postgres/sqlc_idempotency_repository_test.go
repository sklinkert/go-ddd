package postgres

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func TestSqlcIdempotencyRepository_Reserve(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	record := entities.NewIdempotencyRecord("test-key", `{"name": "test product"}`)

	claimed, err := repo.Reserve(ctx, record)

	require.NoError(t, err)
	assert.True(t, claimed)

	// Record is persisted but not completed yet
	foundRecord, err := repo.FindByKey(ctx, "test-key")
	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, record.Id, foundRecord.Id)
	assert.Equal(t, record.Key, foundRecord.Key)
	assert.Equal(t, record.Request, foundRecord.Request)
	assert.Equal(t, "", foundRecord.Response)
	assert.Equal(t, 0, foundRecord.StatusCode)
	assert.False(t, foundRecord.IsCompleted())
	assert.NotEqual(t, uuid.Nil, foundRecord.Id)
	assert.False(t, foundRecord.CreatedAt.IsZero())
}

func TestSqlcIdempotencyRepository_Reserve_DuplicateKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	record1 := entities.NewIdempotencyRecord("duplicate-key", `{"name": "test product 1"}`)
	claimed, err := repo.Reserve(ctx, record1)
	require.NoError(t, err)
	assert.True(t, claimed)

	// Second reserve with the same key must not claim and must not error
	record2 := entities.NewIdempotencyRecord("duplicate-key", `{"name": "test product 2"}`)
	claimed, err = repo.Reserve(ctx, record2)
	require.NoError(t, err)
	assert.False(t, claimed)

	// First record is untouched
	foundRecord, err := repo.FindByKey(ctx, "duplicate-key")
	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, record1.Id, foundRecord.Id)
	assert.Equal(t, record1.Request, foundRecord.Request)
}

func TestSqlcIdempotencyRepository_FindByKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	record := entities.NewIdempotencyRecord("find-test-key", `{"name": "test product"}`)
	claimed, err := repo.Reserve(ctx, record)
	require.NoError(t, err)
	require.True(t, claimed)

	foundRecord, err := repo.FindByKey(ctx, "find-test-key")

	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, record.Id, foundRecord.Id)
	assert.Equal(t, record.Key, foundRecord.Key)
	assert.Equal(t, record.Request, foundRecord.Request)
	assert.Equal(t, record.CreatedAt.Unix(), foundRecord.CreatedAt.Unix())
}

func TestSqlcIdempotencyRepository_FindByKey_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	foundRecord, err := repo.FindByKey(ctx, "non-existent-key")

	require.NoError(t, err) // Should not error, just return nil
	assert.Nil(t, foundRecord)
}

func TestSqlcIdempotencyRepository_SetResponse(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	record := entities.NewIdempotencyRecord("set-response-key", `{"name": "test product"}`)
	claimed, err := repo.Reserve(ctx, record)
	require.NoError(t, err)
	require.True(t, claimed)

	err = repo.SetResponse(ctx, "set-response-key", `{"id": "456", "name": "updated product"}`, 200)
	require.NoError(t, err)

	foundRecord, err := repo.FindByKey(ctx, "set-response-key")
	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, record.Id, foundRecord.Id)
	assert.Equal(t, record.Request, foundRecord.Request)
	assert.Equal(t, `{"id": "456", "name": "updated product"}`, foundRecord.Response)
	assert.Equal(t, 200, foundRecord.StatusCode)
	assert.True(t, foundRecord.IsCompleted())
	assert.Equal(t, record.CreatedAt.Unix(), foundRecord.CreatedAt.Unix()) // CreatedAt should not change
}

func TestSqlcIdempotencyRepository_SetResponse_NonExistentKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Updating a missing key is a no-op, not an error
	err := repo.SetResponse(ctx, "non-existent-key", `{"result": "fail"}`, 404)
	assert.NoError(t, err)

	foundRecord, err := repo.FindByKey(ctx, "non-existent-key")
	require.NoError(t, err)
	assert.Nil(t, foundRecord)
}

func TestSqlcIdempotencyRepository_Delete_ReleasesKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	record := entities.NewIdempotencyRecord("release-key", `{"name": "test product"}`)
	claimed, err := repo.Reserve(ctx, record)
	require.NoError(t, err)
	require.True(t, claimed)

	err = repo.Delete(ctx, "release-key")
	require.NoError(t, err)

	foundRecord, err := repo.FindByKey(ctx, "release-key")
	require.NoError(t, err)
	assert.Nil(t, foundRecord)

	// After release the key can be reserved again
	retryRecord := entities.NewIdempotencyRecord("release-key", `{"name": "retry"}`)
	claimed, err = repo.Reserve(ctx, retryRecord)
	require.NoError(t, err)
	assert.True(t, claimed)
}

func TestSqlcIdempotencyRepository_Delete_NonExistentKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	err := repo.Delete(ctx, "non-existent-key")
	assert.NoError(t, err)
}

func TestSqlcIdempotencyRepository_Reserve_LargeData(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	largeData := make([]byte, 5000)
	for i := range largeData {
		largeData[i] = 'A'
	}
	largeRequest := `{"data": "` + string(largeData) + `"}`
	largeResponse := `{"result": "` + string(largeData) + `"}`

	record := entities.NewIdempotencyRecord("large-data-key", largeRequest)
	claimed, err := repo.Reserve(ctx, record)
	require.NoError(t, err)
	require.True(t, claimed)

	err = repo.SetResponse(ctx, "large-data-key", largeResponse, 200)
	require.NoError(t, err)

	foundRecord, err := repo.FindByKey(ctx, "large-data-key")
	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, largeRequest, foundRecord.Request)
	assert.Equal(t, largeResponse, foundRecord.Response)
	assert.Equal(t, 200, foundRecord.StatusCode)
}

func TestSqlcIdempotencyRepository_Integration_Workflow(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	key := "workflow-test-key"
	requestData := `{"product": "test", "price_cents": 9999, "currency": "USD"}`

	// Step 1: Key does not exist yet
	existingRecord, err := repo.FindByKey(ctx, key)
	require.NoError(t, err)
	assert.Nil(t, existingRecord)

	// Step 2: Reserve the key (request processing started)
	record := entities.NewIdempotencyRecord(key, requestData)
	claimed, err := repo.Reserve(ctx, record)
	require.NoError(t, err)
	require.True(t, claimed)

	// Step 3: While in flight, a concurrent reserve loses the race
	concurrent := entities.NewIdempotencyRecord(key, requestData)
	claimed, err = repo.Reserve(ctx, concurrent)
	require.NoError(t, err)
	assert.False(t, claimed)

	inFlight, err := repo.FindByKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, inFlight)
	assert.False(t, inFlight.IsCompleted())

	// Step 4: Store the response (processing completed)
	responseData := `{"id": "prod-123", "status": "created"}`
	err = repo.SetResponse(ctx, key, responseData, 201)
	require.NoError(t, err)

	finalRecord, err := repo.FindByKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, finalRecord)
	assert.Equal(t, record.Id, finalRecord.Id)
	assert.Equal(t, requestData, finalRecord.Request)
	assert.Equal(t, responseData, finalRecord.Response)
	assert.Equal(t, 201, finalRecord.StatusCode)
	assert.True(t, finalRecord.IsCompleted())
}

func TestSqlcIdempotencyRepository_StatusCodes(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	testCases := []struct {
		name       string
		statusCode int
	}{
		{"success_200", 200},
		{"created_201", 201},
		{"bad_request_400", 400},
		{"not_found_404", 404},
		{"server_error_500", 500},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			key := tc.name + "-key"
			record := entities.NewIdempotencyRecord(key, `{"test": "data"}`)
			claimed, err := repo.Reserve(ctx, record)
			require.NoError(t, err)
			require.True(t, claimed)

			err = repo.SetResponse(ctx, key, `{"result": "test"}`, tc.statusCode)
			require.NoError(t, err)

			foundRecord, err := repo.FindByKey(ctx, key)
			require.NoError(t, err)
			require.NotNil(t, foundRecord)
			assert.Equal(t, tc.statusCode, foundRecord.StatusCode)
		})
	}
}
