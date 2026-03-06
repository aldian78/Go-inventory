package services

import (
	"errors"
	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/repositories"
	"gorm.io/gorm"
)

type ProductService interface {
	GetProducts(name, sku, customer string) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	AdjustStock(id uint, input AdjustmentInput) (*models.Product, error)
}

// AdjustmentInput defines the structure for stock adjustment requests
type AdjustmentInput struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
}

type productService struct {
	productRepo  repositories.ProductRepository
	stockLogRepo repositories.StockLogRepository
	db           *gorm.DB
}

func NewProductService(productRepo repositories.ProductRepository, stockLogRepo repositories.StockLogRepository, db *gorm.DB) ProductService {
	return &productService{
		productRepo:  productRepo,
		stockLogRepo: stockLogRepo,
		db:           db,
	}
}

func (s *productService) GetProducts(name, sku, customer string) ([]models.Product, error) {
	return s.productRepo.FindAll(name, sku, customer)
}

func (s *productService) CreateProduct(product *models.Product) error {
	return s.productRepo.Create(product)
}

func (s *productService) AdjustStock(id uint, input AdjustmentInput) (*models.Product, error) {
	var product *models.Product
	var err error

	// Start Transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Use repositories with transaction handle
		txProductRepo := s.productRepo.WithTrx(tx)
		txStockLogRepo := s.stockLogRepo.WithTrx(tx)

		product, err = txProductRepo.FindByID(id)
		if err != nil {
			return err
		}

		if input.Type == "" {
			input.Type = "BOTH"
		}

		previousPhysical := product.PhysicalStock
		// previousAvailable := product.AvailableStock // Not used in log currently

		switch input.Type {
		case "PHYSICAL":
			product.PhysicalStock += input.Quantity
		case "AVAILABLE":
			product.AvailableStock += input.Quantity
		case "BOTH":
			product.PhysicalStock += input.Quantity
			product.AvailableStock += input.Quantity
		default:
			return errors.New("invalid adjustment type")
		}

		if product.PhysicalStock < 0 || product.AvailableStock < 0 {
			return errors.New("stock cannot be negative after adjustment")
		}

		if err := txProductRepo.Update(product); err != nil {
			return err
		}

		logEntry := models.StockLog{
			ProductID:       product.ID,
			TransactionType: models.TransactionAdjustment,
			Quantity:        input.Quantity,
			PreviousStock:   previousPhysical,
			CurrentStock:    product.PhysicalStock,
			Notes:           input.Notes + " (" + input.Type + ")",
		}

		if err := txStockLogRepo.Create(&logEntry); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}
