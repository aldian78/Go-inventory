package controllers

import (
	"strconv"

	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/services"
	"github.com/gofiber/fiber/v2"
)

type StockInController struct {
	stockInService services.StockInService
}

func NewStockInController(stockInService services.StockInService) *StockInController {
	return &StockInController{stockInService}
}

func (ctrl *StockInController) GetStockIns(c *fiber.Ctx) error {
	stockIns, err := ctrl.stockInService.GetStockIns()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stockIns)
}

func (ctrl *StockInController) CreateStockIn(c *fiber.Ctx) error {
	stockIn := new(models.StockIn)
	if err := c.BodyParser(stockIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ctrl.stockInService.CreateStockIn(stockIn); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(stockIn)
}

func (ctrl *StockInController) UpdateStockInStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock in ID"})
	}

	type UpdateStatusInput struct {
		Status models.StockInStatus `json:"status"`
	}

	var input UpdateStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	stockIn, err := ctrl.stockInService.UpdateStockInStatus(uint(id), input.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stockIn)
}
