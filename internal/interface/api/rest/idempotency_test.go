package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newContext(t *testing.T, header string) echo.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	if header != "" {
		req.Header.Set(idempotencyHeader, header)
	}
	return echo.New().NewContext(req, httptest.NewRecorder())
}

func TestIdempotencyKey_HeaderTakesPrecedence(t *testing.T) {
	c := newContext(t, "from-header")
	assert.Equal(t, "from-header", idempotencyKey(c, "from-body"))
}

func TestIdempotencyKey_FallsBackToBody(t *testing.T) {
	c := newContext(t, "")
	assert.Equal(t, "from-body", idempotencyKey(c, "from-body"))
}

func TestIdempotencyKey_EmptyWhenNeitherPresent(t *testing.T) {
	c := newContext(t, "")
	assert.Empty(t, idempotencyKey(c, ""))
}
