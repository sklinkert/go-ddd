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

// supportedCurrencies maps each supported currency to its ISO 4217
// minor-unit exponent. Not every currency has two decimals: JPY, KRW,
// CLP and ISK have 0; BHD, JOD, KWD, OMR and TND have 3. When adding a
// currency here, use its real exponent so String() stays correct.
var supportedCurrencies = map[Currency]int{
	EUR: 2,
	USD: 2,
}

// Money is an immutable value object storing an amount in ISO 4217 minor
// units (cents for EUR/USD, yen for JPY, fils for BHD) to avoid
// floating-point rounding errors.
type Money struct {
	minorUnits int64
	currency   Currency
}

func NewMoney(minorUnits int64, currency Currency) (Money, error) {
	if minorUnits < 0 {
		return Money{}, fmt.Errorf("%w: amount must not be negative", ErrValidation)
	}
	if _, ok := supportedCurrencies[currency]; !ok {
		return Money{}, fmt.Errorf("%w: unsupported currency %q", ErrValidation, currency)
	}

	return Money{minorUnits: minorUnits, currency: currency}, nil
}

func (m Money) MinorUnits() int64 {
	return m.minorUnits
}

func (m Money) Currency() Currency {
	return m.currency
}

func (m Money) IsZero() bool {
	return m == Money{}
}

func (m Money) String() string {
	exponent := supportedCurrencies[m.currency]
	if exponent == 0 {
		return fmt.Sprintf("%d %s", m.minorUnits, m.currency)
	}

	divisor := int64(1)
	for range exponent {
		divisor *= 10
	}

	return fmt.Sprintf("%d.%0*d %s", m.minorUnits/divisor, exponent, m.minorUnits%divisor, m.currency)
}

type moneyJSON struct {
	MinorUnits int64    `json:"minor_units"`
	Currency   Currency `json:"currency"`

	// LegacyCents catches payloads from before the minor-units rename.
	// Ignoring it would decode an old value as zero and NewMoney would
	// happily accept that, so its presence is an explicit error instead.
	LegacyCents *int64 `json:"cents,omitempty"`
}

func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(moneyJSON{MinorUnits: m.minorUnits, Currency: m.currency})
}

// UnmarshalJSON goes through NewMoney so a Money can never be deserialized
// into an invalid state.
func (m *Money) UnmarshalJSON(data []byte) error {
	var raw moneyJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw.LegacyCents != nil {
		return fmt.Errorf("%w: legacy field %q is no longer supported, use %q", ErrValidation, "cents", "minor_units")
	}

	money, err := NewMoney(raw.MinorUnits, raw.Currency)
	if err != nil {
		return err
	}

	*m = money
	return nil
}
