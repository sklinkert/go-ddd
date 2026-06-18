package testhelpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	db "github.com/sklinkert/go-ddd/internal/infrastructure/db/sqlc"
)

const (
	dbUser     = "testuser"
	dbPassword = "testpass"
	dbName     = "testdb"
)

// PostgresTestContainer manages a test database container
type PostgresTestContainer struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	Queries   *db.Queries
}

// SetupTestDB creates a new PostgreSQL test container and applies the schema
func SetupTestDB(t *testing.T) *PostgresTestContainer {
	ctx := context.Background()

	// Start PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	require.NoError(t, err, "Failed to start postgres container")

	// Get connection string
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "Failed to get connection string")

	// Connect to database
	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err, "Failed to connect to test database")

	// Apply schema
	err = applySchema(ctx, pool)
	require.NoError(t, err, "Failed to apply database schema")

	queries := db.New(pool)

	return &PostgresTestContainer{
		Container: postgresContainer,
		Pool:      pool,
		Queries:   queries,
	}
}

// Close cleans up the test database and container
func (p *PostgresTestContainer) Close(t *testing.T) {
	ctx := context.Background()

	if p.Pool != nil {
		p.Pool.Close()
	}

	if p.Container != nil {
		err := p.Container.Terminate(ctx)
		require.NoError(t, err, "Failed to terminate container")
	}
}

// applySchema reads and applies the SQL schema file
func applySchema(ctx context.Context, pool *pgxpool.Pool) error {
	// Get current working directory and find project root
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Find project root by looking for go.mod file
	projectRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return fmt.Errorf("could not find project root (go.mod not found)")
		}
		projectRoot = parent
	}

	schemaPath := filepath.Join(projectRoot, "migrations", "000001_initial_schema.up.sql")

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = pool.Exec(ctx, string(schemaBytes))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// TruncateTables cleans all test data from the database tables
func (p *PostgresTestContainer) TruncateTables(t *testing.T) {
	ctx := context.Background()

	// Truncate tables in dependency order (child tables first)
	tables := []string{"products", "idempotency_records", "sellers"}

	for _, table := range tables {
		_, err := p.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err, "Failed to truncate table %s", table)
	}
}
