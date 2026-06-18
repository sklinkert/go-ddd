package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sklinkert/go-ddd/internal/application/services"
	postgres2 "github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// DATABASE_URL may be a libpq keyword DSN or a postgres:// URL; pgx accepts both.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=marketplace password=marketplace dbname=marketplace port=5432 sslmode=disable"
	}

	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}

	// Root context is cancelled on SIGINT/SIGTERM so we can drain gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := postgres2.NewConnection(ctx, dsn)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()

	queries := postgres2.NewQueries(pool)

	productRepo := postgres2.NewSqlcProductRepository(queries)
	sellerRepo := postgres2.NewSqlcSellerRepository(queries)
	idempotencyRepo := postgres2.NewSqlcIdempotencyRepository(queries)

	productService := services.NewProductService(productRepo, sellerRepo, idempotencyRepo)
	sellerService := services.NewSellerService(sellerRepo, idempotencyRepo)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(requestLogger(logger))

	rest.NewProductController(e, productService)
	rest.NewSellerController(e, sellerService)

	// Start the server in the background so we can wait for shutdown signals.
	go func() {
		if err := e.Start(port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", slog.Any("error", err))
			stop()
		}
	}()
	logger.Info("server started", slog.String("addr", port))

	<-ctx.Done()
	logger.Info("shutdown signal received; draining in-flight requests")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", slog.Any("error", err))
	}
}

// requestLogger emits one structured log line per request via slog.
func requestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:    true,
		LogURI:       true,
		LogMethod:    true,
		LogLatency:   true,
		LogError:     true,
		LogRequestID: true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("request_id", v.RequestID),
			}
			level := slog.LevelInfo
			if v.Error != nil {
				level = slog.LevelError
				attrs = append(attrs, slog.Any("error", v.Error))
			}
			logger.LogAttrs(c.Request().Context(), level, "request", attrs...)
			return nil
		},
	})
}
