package testhelpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTestDB(t *testing.T) {
	// Test successful database setup
	testDB := SetupTestDB(t)
	require.NotNil(t, testDB)
	require.NotNil(t, testDB.Container)
	require.NotNil(t, testDB.Pool)
	require.NotNil(t, testDB.Queries)

	// Defer cleanup
	defer testDB.Close(t)

	// Verify database connection is working
	ctx := context.Background()
	err := testDB.Pool.Ping(ctx)
	assert.NoError(t, err)
}

func TestPostgresTestContainer_TruncateTables(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Close(t)

	// Create some test data first using raw SQL
	ctx := context.Background()

	// Insert test seller
	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO sellers (id, name, created_at, updated_at) 
		VALUES (gen_random_uuid(), 'Test Seller', NOW(), NOW())
	`)
	require.NoError(t, err)

	// Insert test product
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO products (id, name, price, seller_id, created_at, updated_at)
		VALUES (gen_random_uuid(), 'Test Product', 99.99, 
				(SELECT id FROM sellers LIMIT 1), NOW(), NOW())
	`)
	require.NoError(t, err)

	// Insert test idempotency record
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO idempotency_records (id, key, request, response, status_code, created_at)
		VALUES (gen_random_uuid(), 'test-key', '{"test": "data"}', '{"result": "success"}', 200, NOW())
	`)
	require.NoError(t, err)

	// Verify data exists
	var count int
	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM sellers").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0)

	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0)

	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM idempotency_records").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0)

	// Truncate all tables
	testDB.TruncateTables(t)

	// Verify all tables are empty
	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM sellers").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM idempotency_records").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestPostgresTestContainer_Close(t *testing.T) {
	testDB := SetupTestDB(t)
	require.NotNil(t, testDB)

	// Verify connection is working before close
	ctx := context.Background()
	err := testDB.Pool.Ping(ctx)
	assert.NoError(t, err)

	// Close should not panic and should properly clean up
	assert.NotPanics(t, func() {
		testDB.Close(t)
	})
}

func TestMultipleTestContainers(t *testing.T) {
	// Test that multiple containers can be created simultaneously
	testDB1 := SetupTestDB(t)
	testDB2 := SetupTestDB(t)

	defer testDB1.Close(t)
	defer testDB2.Close(t)

	// Both should be functional
	ctx := context.Background()
	err := testDB1.Pool.Ping(ctx)
	assert.NoError(t, err)

	err = testDB2.Pool.Ping(ctx)
	assert.NoError(t, err)

	// They should be independent
	assert.NotEqual(t, testDB1.Container, testDB2.Container)
	assert.NotEqual(t, testDB1.Pool, testDB2.Pool)
}
