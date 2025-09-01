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

func TestNumericFromFloat64(t *testing.T) {
	testCases := []float64{
		0.0,
		99.99,
		123.45,
		-50.25,
		0.01,
		999999.99,
	}

	for _, testValue := range testCases {
		t.Run(formatFloat(testValue), func(t *testing.T) {
			pgNumeric := numericFromFloat64(testValue)

			assert.True(t, pgNumeric.Valid)

			// Convert back to verify
			floatValue, err := pgNumeric.Float64Value()
			assert.NoError(t, err)
			assert.True(t, floatValue.Valid)
			assert.InDelta(t, testValue, floatValue.Float64, 0.01) // Allow for precision differences
		})
	}
}

func TestFloat64FromNumeric_Valid(t *testing.T) {
	// Create a valid numeric
	var pgNumeric pgtype.Numeric
	err := pgNumeric.Scan("123.45")
	assert.NoError(t, err)

	result := float64FromNumeric(pgNumeric)

	assert.InDelta(t, 123.45, result, 0.01)
}

func TestFloat64FromNumeric_Invalid(t *testing.T) {
	// Create an invalid numeric
	var pgNumeric pgtype.Numeric
	// Don't scan anything, leaving it invalid

	result := float64FromNumeric(pgNumeric)

	// Should return 0 for invalid numeric
	assert.Equal(t, 0.0, result)
}

func TestRoundTripConversions(t *testing.T) {
	// Test round-trip conversions to ensure data integrity

	// Time round-trip
	originalTime := time.Now().Truncate(time.Microsecond) // Truncate to DB precision
	pgTime := timestamptzFromTime(originalTime)
	resultTime := timeFromTimestamptz(pgTime)
	assert.True(t, originalTime.Equal(resultTime))

	// Float64 round-trip
	originalFloat := 123.45
	pgNumeric := numericFromFloat64(originalFloat)
	resultFloat := float64FromNumeric(pgNumeric)
	assert.InDelta(t, originalFloat, resultFloat, 0.01)
}

func TestEdgeCaseNumbers(t *testing.T) {
	testCases := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"small positive", 0.01},
		{"small negative", -0.01},
		{"large positive", 99999999.99},
		{"large negative", -99999999.99},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pgNumeric := numericFromFloat64(tc.value)
			result := float64FromNumeric(pgNumeric)
			assert.InDelta(t, tc.value, result, 0.01)
		})
	}
}

// Helper function to format float for test names
func formatFloat(f float64) string {
	if f == 0.0 {
		return "zero"
	} else if f > 0 {
		return "positive"
	}
	return "negative"
}
