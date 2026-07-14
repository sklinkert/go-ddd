package entities

import (
	"encoding/json"
	"fmt"
)

// Currency is the ISO 4217 code of a Money value.
type Currency string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

var supportedCurrencies = map[Currency]struct{}{
	EUR: {},
	USD: {},
}

// Money is an immutable value object storing an amount in minor units
// (cents) to avoid floating-point rounding errors.
type Money struct {
	cents    int64
	currency Currency
}

func NewMoney(cents int64, currency Currency) (Money, error) {
	if cents < 0 {
		return Money{}, fmt.Errorf("%w: amount must not be negative", ErrValidation)
	}
	if _, ok := supportedCurrencies[currency]; !ok {
		return Money{}, fmt.Errorf("%w: unsupported currency %q", ErrValidation, currency)
	}

	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) IsZero() bool {
	return m == Money{}
}

func (m Money) String() string {
	return fmt.Sprintf("%d.%02d %s", m.cents/100, m.cents%100, m.currency)
}

type moneyJSON struct {
	Cents    int64    `json:"cents"`
	Currency Currency `json:"currency"`
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(moneyJSON{Cents: m.cents, Currency: m.currency})
}

// UnmarshalJSON goes through NewMoney so a Money can never be deserialized
// into an invalid state.
func (m *Money) UnmarshalJSON(data []byte) error {
	var raw moneyJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	money, err := NewMoney(raw.Cents, raw.Currency)
	if err != nil {
		return err
	}

	*m = money
	return nil
}
