package entities

import "errors"

// Sentinel errors that the interface layer can translate into HTTP status
// codes (e.g. 404) instead of a generic 500.
var (
	ErrProductNotFound = errors.New("product not found")
	ErrSellerNotFound  = errors.New("seller not found")
	// ErrValidation wraps all domain invariant violations; check with
	// errors.Is to translate into a 400.
	ErrValidation = errors.New("validation failed")
)
