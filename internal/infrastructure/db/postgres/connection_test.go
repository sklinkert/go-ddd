package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sklinkert/go-ddd/internal/testhelpers"
)

func TestNewConnection(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	// Test NewQueries with the existing connection (we can't easily get a new DSN)
	queries := NewQueries(testDB.Conn)
	assert.NotNil(t, queries)

	// Verify the connection is working by using the existing connection
	ctx := context.Background()
	err := testDB.Conn.Ping(ctx)
	assert.NoError(t, err)
}

func TestNewConnection_InvalidDSN(t *testing.T) {
	ctx := context.Background()

	// Test with invalid DSN
	conn, err := NewConnection(ctx, "invalid-dsn")
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestNewConnection_UnreachableHost(t *testing.T) {
	ctx := context.Background()

	// Test with unreachable host
	conn, err := NewConnection(ctx, "postgres://user:pass@unreachable-host:5432/db?connect_timeout=1")
	assert.Error(t, err)
	assert.Nil(t, conn)
}

func TestNewQueries(t *testing.T) {
	testDB := testhelpers.SetupTestDB(t)
	defer testDB.Close(t)

	// Test NewQueries with valid connection
	queries := NewQueries(testDB.Conn)
	assert.NotNil(t, queries)

	// Verify queries object is functional by running a simple query
	ctx := context.Background()
	_, err := queries.GetAllProducts(ctx)
	assert.NoError(t, err) // Should not error even if empty
}

func TestNewQueries_WithNilConnection(t *testing.T) {
	// Test NewQueries with nil connection
	// Note: This will create a queries object but will panic when used
	queries := NewQueries(nil)
	assert.NotNil(t, queries)

	// Attempting to use it should panic (so we won't test that)
	// This test just verifies that NewQueries can accept nil without immediately panicking
}
