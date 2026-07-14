package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestTimestamptzFromTime(t *testing.T) {
	// Test with current time
	now := time.Now()
	pgTime := timestamptzFromTime(now)

	assert.True(t, pgTime.Valid)
	// Convert back to compare (allowing for microsecond precision differences)
	assert.True(t, pgTime.Time.Sub(now).Abs() < time.Microsecond)
}

func TestTimestamptzFromTime_ZeroTime(t *testing.T) {
	// Test with zero time
	zeroTime := time.Time{}
	pgTime := timestamptzFromTime(zeroTime)

	// Should still be valid but represent zero time
	assert.True(t, pgTime.Valid)
	assert.True(t, pgTime.Time.IsZero())
}

func TestTimeFromTimestamptz_Valid(t *testing.T) {
	// Create a valid timestamptz
	now := time.Now()
	var pgTime pgtype.Timestamptz
	err := pgTime.Scan(now)
	assert.NoError(t, err)

	// Convert back to time
	result := timeFromTimestamptz(pgTime)

	// Should be approximately equal (allowing for precision loss)
	assert.True(t, result.Sub(now).Abs() < time.Microsecond)
}

func TestTimeFromTimestamptz_Invalid(t *testing.T) {
	// Create an invalid timestamptz
	var pgTime pgtype.Timestamptz
	// Don't scan anything, leaving it invalid

	result := timeFromTimestamptz(pgTime)

	// Should return zero time for invalid input
	assert.True(t, result.IsZero())
}

func TestTimestamptz_RoundTrip(t *testing.T) {
	originalTime := time.Now().Truncate(time.Microsecond) // Truncate to DB precision
	pgTime := timestamptzFromTime(originalTime)
	resultTime := timeFromTimestamptz(pgTime)
	assert.True(t, originalTime.Equal(resultTime))
}
