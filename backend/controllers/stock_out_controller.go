package controllers

import (
	"strconv"

	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/services"
	"github.com/gofiber/fiber/v2"
)

type StockOutController struct {
	stockOutService services.StockOutService
}

func NewStockOutController(stockOutService services.StockOutService) *StockOutController {
	return &StockOutController{stockOutService}
}

func (ctrl *StockOutController) GetStockOuts(c *fiber.Ctx) error {
	stockOuts, err := ctrl.stockOutService.GetStockOuts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stockOuts)
}

func (ctrl *StockOutController) CreateStockOut(c *fiber.Ctx) error {
	stockOut := new(models.StockOut)
	if err := c.BodyParser(stockOut); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	createdStockOut, err := ctrl.stockOutService.CreateStockOut(stockOut)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(createdStockOut)
}

func (ctrl *StockOutController) UpdateStockOutStatus(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock out ID"})
	}

	type UpdateStatusInput struct {
		Status models.StockOutStatus `json:"status"`
	}

	var input UpdateStatusInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	stockOut, err := ctrl.stockOutService.UpdateStockOutStatus(uint(id), input.Status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(stockOut)
}
