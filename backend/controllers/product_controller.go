package controllers

import (
	"strconv"

	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/services"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{productService}
}

func (ctrl *ProductController) GetProducts(c *fiber.Ctx) error {
	name := c.Query("name")
	sku := c.Query("sku")
	customer := c.Query("customer")

	products, err := ctrl.productService.GetProducts(name, sku, customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}

func (ctrl *ProductController) CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ctrl.productService.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

func (ctrl *ProductController) AdjustStock(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	var input services.AdjustmentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	product, err := ctrl.productService.AdjustStock(uint(id), input)
	if err != nil {
		// You might want to check for specific error types here to return different status codes
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(product)
}
