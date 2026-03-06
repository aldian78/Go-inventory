package services

import (
	"errors"
	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/repositories"
	"gorm.io/gorm"
)

type StockInService interface {
	GetStockIns() ([]models.StockIn, error)
	CreateStockIn(stockIn *models.StockIn) error
	UpdateStockInStatus(id uint, status models.StockInStatus) (*models.StockIn, error)
}

type stockInService struct {
	stockInRepo  repositories.StockInRepository
	productRepo  repositories.ProductRepository
	stockLogRepo repositories.StockLogRepository
	db           *gorm.DB
}

func NewStockInService(
	stockInRepo repositories.StockInRepository,
	productRepo repositories.ProductRepository,
	stockLogRepo repositories.StockLogRepository,
	db *gorm.DB,
) StockInService {
	return &stockInService{
		stockInRepo:  stockInRepo,
		productRepo:  productRepo,
		stockLogRepo: stockLogRepo,
		db:           db,
	}
}

func (s *stockInService) GetStockIns() ([]models.StockIn, error) {
	return s.stockInRepo.FindAll()
}

func (s *stockInService) CreateStockIn(stockIn *models.StockIn) error {
	return s.stockInRepo.Create(stockIn)
}

func (s *stockInService) UpdateStockInStatus(id uint, status models.StockInStatus) (*models.StockIn, error) {
	var stockIn *models.StockIn
	var err error

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txStockInRepo := s.stockInRepo.WithTrx(tx)
		txProductRepo := s.productRepo.WithTrx(tx)
		txStockLogRepo := s.stockLogRepo.WithTrx(tx)

		stockIn, err = txStockInRepo.FindByID(id)
		if err != nil {
			return err
		}

		if stockIn.Status == models.StockInDone || stockIn.Status == models.StockInCancelled {
			return errors.New("transaction is already final")
		}

		originalStatus := stockIn.Status
		stockIn.Status = status

		if status == models.StockInDone {
			if originalStatus != models.StockInInProgress {
				return errors.New("can only complete transactions that are IN_PROGRESS")
			}

			product, err := txProductRepo.FindByID(stockIn.ProductID)
			if err != nil {
				return err
			}

			previousStock := product.PhysicalStock
			product.PhysicalStock += stockIn.Quantity
			product.AvailableStock += stockIn.Quantity

			if err := txProductRepo.Update(product); err != nil {
				return err
			}

			logEntry := models.StockLog{
				ProductID:       stockIn.ProductID,
				TransactionID:   stockIn.ID,
				TransactionType: models.TransactionStockIn,
				Quantity:        stockIn.Quantity,
				PreviousStock:   previousStock,
				CurrentStock:    product.PhysicalStock,
				Notes:           "Stock In Completed",
			}
			if err := txStockLogRepo.Create(&logEntry); err != nil {
				return err
			}
		}

		if err := txStockInRepo.Update(stockIn); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return stockIn, nil
}
