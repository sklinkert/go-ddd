package rest

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type HealthController struct {
	pool *pgxpool.Pool
}

// NewHealthController exposes liveness and readiness probes. Liveness only
// proves the process responds; readiness also checks the database, so
// orchestrators stop routing traffic when Postgres is gone.
func NewHealthController(e *echo.Echo, pool *pgxpool.Pool) *HealthController {
	controller := &HealthController{pool: pool}

	e.GET("/healthz", controller.Liveness)
	e.GET("/readyz", controller.Readiness)

	return controller
}

func (hc *HealthController) Liveness(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (hc *HealthController) Readiness(c echo.Context) error {
	if err := hc.pool.Ping(c.Request().Context()); err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
			"reason": "database unreachable",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
