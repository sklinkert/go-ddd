package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sklinkert/go-ddd/internal/domain/entities"
	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func TestSqlcIdempotencyRepository_Create(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create an idempotency record
	record := entities.NewIdempotencyRecord("test-key", `{"name": "test product"}`)
	record.SetResponse(`{"id": "123", "name": "test product"}`, 201)

	createdRecord, err := repo.Create(ctx, record)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, createdRecord)
	assert.Equal(t, record.Key, createdRecord.Key)
	assert.Equal(t, record.Request, createdRecord.Request)
	assert.Equal(t, record.Response, createdRecord.Response)
	assert.Equal(t, record.StatusCode, createdRecord.StatusCode)
	assert.NotEqual(t, uuid.Nil, createdRecord.ID)
	assert.False(t, createdRecord.CreatedAt.IsZero())
}

func TestSqlcIdempotencyRepository_FindByKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create test data
	record := entities.NewIdempotencyRecord("find-test-key", `{"name": "test product"}`)
	record.SetResponse(`{"id": "123", "name": "test product"}`, 201)

	createdRecord, err := repo.Create(ctx, record)
	require.NoError(t, err)

	// Test finding by key
	foundRecord, err := repo.FindByKey(ctx, "find-test-key")

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, foundRecord)
	assert.Equal(t, createdRecord.ID, foundRecord.ID)
	assert.Equal(t, createdRecord.Key, foundRecord.Key)
	assert.Equal(t, createdRecord.Request, foundRecord.Request)
	assert.Equal(t, createdRecord.Response, foundRecord.Response)
	assert.Equal(t, createdRecord.StatusCode, foundRecord.StatusCode)
	assert.Equal(t, createdRecord.CreatedAt.Unix(), foundRecord.CreatedAt.Unix())
}

func TestSqlcIdempotencyRepository_FindByKey_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Test finding non-existent key
	foundRecord, err := repo.FindByKey(ctx, "non-existent-key")

	// Assertions
	require.NoError(t, err) // Should not error, just return nil
	assert.Nil(t, foundRecord)
}

func TestSqlcIdempotencyRepository_Update(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create test data
	record := entities.NewIdempotencyRecord("update-test-key", `{"name": "test product"}`)
	// Initially no response set

	createdRecord, err := repo.Create(ctx, record)
	require.NoError(t, err)

	// Update the record with response
	createdRecord.SetResponse(`{"id": "456", "name": "updated product"}`, 200)

	updatedRecord, err := repo.Update(ctx, createdRecord)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, updatedRecord)
	assert.Equal(t, createdRecord.ID, updatedRecord.ID)
	assert.Equal(t, createdRecord.Key, updatedRecord.Key)
	assert.Equal(t, createdRecord.Request, updatedRecord.Request)
	assert.Equal(t, `{"id": "456", "name": "updated product"}`, updatedRecord.Response)
	assert.Equal(t, 200, updatedRecord.StatusCode)
	assert.Equal(t, createdRecord.CreatedAt.Unix(), updatedRecord.CreatedAt.Unix()) // CreatedAt should not change
}

func TestSqlcIdempotencyRepository_Update_NotFound(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Try to update non-existent record
	nonExistentRecord := &entities.IdempotencyRecord{
		ID:         uuid.New(),
		Key:        "non-existent-key",
		Request:    `{"test": "data"}`,
		Response:   `{"result": "fail"}`,
		StatusCode: 404,
		CreatedAt:  time.Now(),
	}

	updatedRecord, err := repo.Update(ctx, nonExistentRecord)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, updatedRecord)
}

func TestSqlcIdempotencyRepository_Create_DuplicateKey(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create first record
	record1 := entities.NewIdempotencyRecord("duplicate-key", `{"name": "test product 1"}`)
	record1.SetResponse(`{"id": "123"}`, 201)

	_, err := repo.Create(ctx, record1)
	require.NoError(t, err)

	// Try to create second record with same key
	record2 := entities.NewIdempotencyRecord("duplicate-key", `{"name": "test product 2"}`)
	record2.SetResponse(`{"id": "456"}`, 201)

	createdRecord2, err := repo.Create(ctx, record2)

	// Should fail due to unique constraint on key
	assert.Error(t, err)
	assert.Nil(t, createdRecord2)
}

func TestSqlcIdempotencyRepository_Create_EmptyResponse(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create record without setting response (empty response, zero status code)
	record := entities.NewIdempotencyRecord("empty-response-key", `{"name": "test product"}`)
	// Note: not calling SetResponse, so Response is empty string and StatusCode is 0

	createdRecord, err := repo.Create(ctx, record)

	// Should succeed with empty response
	require.NoError(t, err)
	require.NotNil(t, createdRecord)
	assert.Equal(t, "empty-response-key", createdRecord.Key)
	assert.Equal(t, `{"name": "test product"}`, createdRecord.Request)
	assert.Equal(t, "", createdRecord.Response)
	assert.Equal(t, 0, createdRecord.StatusCode)
}

func TestSqlcIdempotencyRepository_Create_LargeData(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	// Create record with large request/response data
	largeData := make([]byte, 5000)
	for i := range largeData {
		largeData[i] = 'A'
	}
	largeRequest := `{"data": "` + string(largeData) + `"}`
	largeResponse := `{"result": "` + string(largeData) + `"}`

	record := entities.NewIdempotencyRecord("large-data-key", largeRequest)
	record.SetResponse(largeResponse, 200)

	createdRecord, err := repo.Create(ctx, record)

	// Should succeed as TEXT fields can handle large data
	require.NoError(t, err)
	require.NotNil(t, createdRecord)
	assert.Equal(t, "large-data-key", createdRecord.Key)
	assert.Equal(t, largeRequest, createdRecord.Request)
	assert.Equal(t, largeResponse, createdRecord.Response)
	assert.Equal(t, 200, createdRecord.StatusCode)
}

func TestSqlcIdempotencyRepository_Integration_Workflow(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	repo := NewSqlcIdempotencyRepository(testDB.Queries)
	ctx := context.Background()

	key := "workflow-test-key"
	requestData := `{"product": "test", "price": 99.99}`

	// Step 1: Check if key exists (should not exist)
	existingRecord, err := repo.FindByKey(ctx, key)
	require.NoError(t, err)
	assert.Nil(t, existingRecord)

	// Step 2: Create new record without response (request processing started)
	initialRecord := entities.NewIdempotencyRecord(key, requestData)

	createdRecord, err := repo.Create(ctx, initialRecord)
	require.NoError(t, err)
	require.NotNil(t, createdRecord)
	assert.Equal(t, key, createdRecord.Key)
	assert.Equal(t, requestData, createdRecord.Request)
	assert.Equal(t, "", createdRecord.Response) // No response yet
	assert.Equal(t, 0, createdRecord.StatusCode)

	// Step 3: Update record with response (processing completed)
	responseData := `{"id": "prod-123", "status": "created"}`
	createdRecord.SetResponse(responseData, 201)

	updatedRecord, err := repo.Update(ctx, createdRecord)
	require.NoError(t, err)
	require.NotNil(t, updatedRecord)
	assert.Equal(t, key, updatedRecord.Key)
	assert.Equal(t, requestData, updatedRecord.Request)
	assert.Equal(t, responseData, updatedRecord.Response)
	assert.Equal(t, 201, updatedRecord.StatusCode)

	// Step 4: Retrieve final record by key
	finalRecord, err := repo.FindByKey(ctx, key)
	require.NoError(t, err)
	require.NotNil(t, finalRecord)
	assert.Equal(t, updatedRecord.ID, finalRecord.ID)
	assert.Equal(t, updatedRecord.Key, finalRecord.Key)
	assert.Equal(t, updatedRecord.Request, finalRecord.Request)
	assert.Equal(t, updatedRecord.Response, finalRecord.Response)
	assert.Equal(t, updatedRecord.StatusCode, finalRecord.StatusCode)
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
		{"zero_status", 0},
		{"negative_status", -1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			record := entities.NewIdempotencyRecord(tc.name+"-key", `{"test": "data"}`)
			record.SetResponse(`{"result": "test"}`, tc.statusCode)

			createdRecord, err := repo.Create(ctx, record)

			require.NoError(t, err)
			require.NotNil(t, createdRecord)
			assert.Equal(t, tc.statusCode, createdRecord.StatusCode)

			// Verify by retrieving
			foundRecord, err := repo.FindByKey(ctx, tc.name+"-key")
			require.NoError(t, err)
			require.NotNil(t, foundRecord)
			assert.Equal(t, tc.statusCode, foundRecord.StatusCode)
		})
	}
}
