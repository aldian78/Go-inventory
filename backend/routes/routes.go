package routes

import (
	"github.com/go-inventory/backend/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	productController *controllers.ProductController,
	stockInController *controllers.StockInController,
	stockOutController *controllers.StockOutController,
	reportController *controllers.ReportController,
) {
	api := app.Group("/api")

	// Product Routes
	api.Get("/products", productController.GetProducts)
	api.Post("/products", productController.CreateProduct)
	api.Put("/products/:id/adjust", productController.AdjustStock)

	// Stock In Routes
	api.Get("/stock-in", stockInController.GetStockIns)
	api.Post("/stock-in", stockInController.CreateStockIn)
	api.Put("/stock-in/:id/status", stockInController.UpdateStockInStatus)

	// Stock Out Routes
	api.Get("/stock-out", stockOutController.GetStockOuts)
	api.Post("/stock-out", stockOutController.CreateStockOut)
	api.Put("/stock-out/:id/status", stockOutController.UpdateStockOutStatus)

	// Report Routes
	api.Get("/reports/transactions", reportController.GetTransactionReport)
}
