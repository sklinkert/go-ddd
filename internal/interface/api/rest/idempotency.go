package rest

import "github.com/labstack/echo/v4"

// idempotencyHeader is the conventional header clients use to make a mutating
// request safe to retry. See https://datatracker.ietf.org/doc/draft-ietf-httpapi-idempotency-key-header/
const idempotencyHeader = "Idempotency-Key"

// idempotencyKey resolves the idempotency key for a request. The Idempotency-Key
// header is the preferred transport; bodyKey (the request's optional
// "idempotency_key" field) is a backward-compatible fallback for callers that
// send the key in the JSON body. The header wins when both are present.
func idempotencyKey(c echo.Context, bodyKey string) string {
	if header := c.Request().Header.Get(idempotencyHeader); header != "" {
		return header
	}
	return bodyKey
}
