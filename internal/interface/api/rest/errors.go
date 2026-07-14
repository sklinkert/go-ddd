package rest

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/domain/entities"
)

// writeCommandError maps well-known service errors to HTTP status codes so
// clients get 404/409 instead of a generic 500.
func writeCommandError(c echo.Context, err error, fallback string) error {
	switch {
	case errors.Is(err, entities.ErrProductNotFound), errors.Is(err, entities.ErrSellerNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, entities.ErrValidation):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, services.ErrRequestInFlight):
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, services.ErrIdempotencyKeyReuse):
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fallback})
	}
}
