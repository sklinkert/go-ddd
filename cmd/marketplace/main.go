package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	postgres2 "github.com/sklinkert/go-ddd/internal/infrastructure/db/postgres"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"log"
)

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	port := ":8080"

	ctx := context.Background()
	conn, err := postgres2.NewConnection(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	queries := postgres2.NewQueries(conn)

	productRepo := postgres2.NewSqlcProductRepository(queries)
	sellerRepo := postgres2.NewSqlcSellerRepository(queries)
	idempotencyRepo := postgres2.NewSqlcIdempotencyRepository(queries)

	productService := services.NewProductService(productRepo, sellerRepo, idempotencyRepo)
	sellerService := services.NewSellerService(sellerRepo, idempotencyRepo)

	e := echo.New()
	rest.NewProductController(e, productService)
	rest.NewSellerController(e, sellerService)

	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
