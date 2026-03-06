package main

import (
	"log"

	"github.com/go-inventory/backend/config"
	"github.com/go-inventory/backend/controllers"
	"github.com/go-inventory/backend/repositories"
	"github.com/go-inventory/backend/routes"
	"github.com/go-inventory/backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(cors.New())

	// Database
	config.ConnectDB()
	db := config.DB

	// Repositories
	productRepo := repositories.NewProductRepository(db)
	stockInRepo := repositories.NewStockInRepository(db)
	stockOutRepo := repositories.NewStockOutRepository(db)
	stockLogRepo := repositories.NewStockLogRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	// Services
	productService := services.NewProductService(productRepo, stockLogRepo, db)
	stockInService := services.NewStockInService(stockInRepo, productRepo, stockLogRepo, db)
	stockOutService := services.NewStockOutService(stockOutRepo, productRepo, stockLogRepo, db)
	reportService := services.NewReportService(reportRepo)

	// Controllers
	productController := controllers.NewProductController(productService)
	stockInController := controllers.NewStockInController(stockInService)
	stockOutController := controllers.NewStockOutController(stockOutService)
	reportController := controllers.NewReportController(reportService)

	// Routes
	routes.SetupRoutes(app, productController, stockInController, stockOutController, reportController)

	log.Fatal(app.Listen(":8080"))
}
