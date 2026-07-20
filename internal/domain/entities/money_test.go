package entities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMoney(t *testing.T) {
	money, err := NewMoney(1234, USD)
	require.NoError(t, err)

	assert.Equal(t, int64(1234), money.MinorUnits())
	assert.Equal(t, USD, money.Currency())
	assert.False(t, money.IsZero())
}

func TestNewMoney_ZeroAmountAllowed(t *testing.T) {
	money, err := NewMoney(0, EUR)
	require.NoError(t, err)

	assert.Equal(t, int64(0), money.MinorUnits())
	assert.Equal(t, EUR, money.Currency())
}

func TestNewMoney_NegativeAmount(t *testing.T) {
	_, err := NewMoney(-1, USD)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be negative")
}

func TestNewMoney_UnsupportedCurrency(t *testing.T) {
	testCases := []Currency{"", "GBP", "usd"}

	for _, currency := range testCases {
		t.Run(string(currency), func(t *testing.T) {
			_, err := NewMoney(100, currency)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "unsupported currency")
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	assert.True(t, Money{}.IsZero())

	money, err := NewMoney(1, USD)
	require.NoError(t, err)
	assert.False(t, money.IsZero())

	// A zero amount with a currency is not the zero value
	zeroUsd, err := NewMoney(0, USD)
	require.NoError(t, err)
	assert.False(t, zeroUsd.IsZero())
}

func TestMoney_String(t *testing.T) {
	testCases := []struct {
		minorUnits int64
		currency   Currency
		expected   string
	}{
		{1234, USD, "12.34 USD"},
		{5, EUR, "0.05 EUR"},
		{100, EUR, "1.00 EUR"},
		{999999, USD, "9999.99 USD"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			money, err := NewMoney(tc.minorUnits, tc.currency)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, money.String())
		})
	}
}

// TestMoney_String_ExponentAware pins the formatting contract for
// non-exponent-2 currencies without adding them to the whitelist: the
// divisor comes from the exponent table, not a hardcoded 100.
func TestMoney_String_ExponentAware(t *testing.T) {
	const (
		testJPY Currency = "JPY" // exponent 0
		testBHD Currency = "BHD" // exponent 3
	)
	supportedCurrencies[testJPY] = 0
	supportedCurrencies[testBHD] = 3
	defer delete(supportedCurrencies, testJPY)
	defer delete(supportedCurrencies, testBHD)

	testCases := []struct {
		minorUnits int64
		currency   Currency
		expected   string
	}{
		{5000, testJPY, "5000 JPY"},
		{0, testJPY, "0 JPY"},
		{12345, testBHD, "12.345 BHD"},
		{5, testBHD, "0.005 BHD"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			money, err := NewMoney(tc.minorUnits, tc.currency)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, money.String())
		})
	}
}

func TestMoney_JSONRoundTrip(t *testing.T) {
	original, err := NewMoney(1234, USD)
	require.NoError(t, err)

	data, err := json.Marshal(original)
	require.NoError(t, err)
	assert.JSONEq(t, `{"minor_units":1234,"currency":"USD"}`, string(data))

	var decoded Money
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}

func TestMoney_UnmarshalJSON_InvalidState(t *testing.T) {
	testCases := []struct {
		name string
		json string
	}{
		{"negative amount", `{"minor_units":-100,"currency":"USD"}`},
		{"unsupported currency", `{"minor_units":100,"currency":"GBP"}`},
		{"missing currency", `{"minor_units":100}`},
		{"malformed json", `{"minor_units":`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var money Money
			err := json.Unmarshal([]byte(tc.json), &money)
			assert.Error(t, err)
		})
	}
}

// A pre-rename payload must fail loudly: silently ignoring "cents" would
// decode as a zero amount, which NewMoney accepts.
func TestMoney_UnmarshalJSON_RejectsLegacyCentsField(t *testing.T) {
	var money Money
	err := json.Unmarshal([]byte(`{"cents":1234,"currency":"USD"}`), &money)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrValidation)
	assert.Contains(t, err.Error(), "minor_units")
}
