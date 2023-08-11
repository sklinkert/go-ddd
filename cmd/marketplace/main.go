package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/sklinkert/go-ddd/internal/application/services"
	"github.com/sklinkert/go-ddd/internal/infrastructure/db"
	"github.com/sklinkert/go-ddd/internal/interface/api/rest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	port := ":8080"

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	//gormDB.AutoMigrate()

	productRepo := db.NewGormProductRepository(gormDB)
	sellerRepo := db.NewGormSellerRepository(gormDB)

	productService := services.NewProductService(productRepo)
	sellerService := services.NewSellerService(sellerRepo)

	e := echo.New()
	rest.NewProductController(e, productService)
	rest.NewSellerController(e, sellerService)

	if err := e.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
