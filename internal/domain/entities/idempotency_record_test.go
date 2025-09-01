package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewIdempotencyRecord(t *testing.T) {
	key := "test-key-123"
	request := `{"product": "test", "price": 99.99}`

	record := NewIdempotencyRecord(key, request)

	// Assertions
	assert.NotEqual(t, uuid.Nil, record.Id)
	assert.Equal(t, key, record.Key)
	assert.Equal(t, request, record.Request)
	assert.Equal(t, "", record.Response)  // Should be empty initially
	assert.Equal(t, 0, record.StatusCode) // Should be 0 initially
	assert.False(t, record.CreatedAt.IsZero())
}

func TestIdempotencyRecord_SetResponse(t *testing.T) {
	record := NewIdempotencyRecord("test-key", `{"test": "request"}`)

	// Initially no response
	assert.Equal(t, "", record.Response)
	assert.Equal(t, 0, record.StatusCode)

	// Set response
	response := `{"id": "123", "status": "created"}`
	statusCode := 201

	record.SetResponse(response, statusCode)

	// Verify response is set
	assert.Equal(t, response, record.Response)
	assert.Equal(t, statusCode, record.StatusCode)
}

func TestIdempotencyRecord_SetResponse_Multiple(t *testing.T) {
	record := NewIdempotencyRecord("test-key", `{"test": "request"}`)

	// Set initial response
	record.SetResponse(`{"status": "processing"}`, 202)
	assert.Equal(t, `{"status": "processing"}`, record.Response)
	assert.Equal(t, 202, record.StatusCode)

	// Update response (simulating completion)
	record.SetResponse(`{"id": "123", "status": "completed"}`, 200)
	assert.Equal(t, `{"id": "123", "status": "completed"}`, record.Response)
	assert.Equal(t, 200, record.StatusCode)
}

func TestIdempotencyRecord_SetResponse_ErrorScenarios(t *testing.T) {
	testCases := []struct {
		name       string
		response   string
		statusCode int
	}{
		{"client error", `{"error": "Bad Request"}`, 400},
		{"server error", `{"error": "Internal Server Error"}`, 500},
		{"empty response", "", 204},
		{"negative status", `{"error": "test"}`, -1},
		{"large response", string(make([]byte, 1000)), 200},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset record for each test
			testRecord := NewIdempotencyRecord(tc.name+"-key", `{"test": "request"}`)

			testRecord.SetResponse(tc.response, tc.statusCode)

			assert.Equal(t, tc.response, testRecord.Response)
			assert.Equal(t, tc.statusCode, testRecord.StatusCode)
		})
	}
}

func TestIdempotencyRecord_ImmutableFields(t *testing.T) {
	key := "immutable-test-key"
	request := `{"test": "immutable"}`

	record := NewIdempotencyRecord(key, request)
	originalID := record.Id
	originalKey := record.Key
	originalRequest := record.Request
	originalCreatedAt := record.CreatedAt

	// Set response multiple times
	record.SetResponse(`{"status": "first"}`, 200)
	record.SetResponse(`{"status": "second"}`, 201)

	// Verify immutable fields haven't changed
	assert.Equal(t, originalID, record.Id)
	assert.Equal(t, originalKey, record.Key)
	assert.Equal(t, originalRequest, record.Request)
	assert.Equal(t, originalCreatedAt, record.CreatedAt)

	// Verify mutable fields have changed
	assert.Equal(t, `{"status": "second"}`, record.Response)
	assert.Equal(t, 201, record.StatusCode)
}

func TestIdempotencyRecord_UniqueIDs(t *testing.T) {
	// Create multiple records to ensure IDs are unique
	records := make([]*IdempotencyRecord, 100)
	ids := make(map[uuid.UUID]bool)

	for i := 0; i < 100; i++ {
		record := NewIdempotencyRecord("key-"+string(rune(i)), `{"test": "data"}`)
		records[i] = record

		// Check for duplicate IDs
		assert.False(t, ids[record.Id], "Duplicate ID found: %s", record.Id)
		ids[record.Id] = true
	}

	// All IDs should be unique
	assert.Len(t, ids, 100)
}

func TestIdempotencyRecord_CreatedAtConsistency(t *testing.T) {
	before := time.Now()
	record := NewIdempotencyRecord("time-test-key", `{"test": "time"}`)
	after := time.Now()

	// CreatedAt should be between before and after
	assert.True(t, record.CreatedAt.After(before) || record.CreatedAt.Equal(before))
	assert.True(t, record.CreatedAt.Before(after) || record.CreatedAt.Equal(after))
}
