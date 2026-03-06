package controllers

import (
	"time"

	"github.com/go-inventory/backend/services"
	"github.com/gofiber/fiber/v2"
)

type ReportController struct {
	reportService services.ReportService
}

func NewReportController(reportService services.ReportService) *ReportController {
	return &ReportController{reportService}
}

func (ctrl *ReportController) GetTransactionReport(c *fiber.Ctx) error {
	transType := c.Query("type") // Expected: "stock_in", "stock_out", or "adjustment"
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if transType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'type' is required"})
	}

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		}
	}

	reportData, err := ctrl.reportService.GetTransactionReport(transType, startDate, endDate)
	if err != nil {
		// Check for specific service errors if needed
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reportData)
}
