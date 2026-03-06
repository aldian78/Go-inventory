package services

import (
	"errors"
	"github.com/go-inventory/backend/models"
	"github.com/go-inventory/backend/repositories"
	"gorm.io/gorm"
)

type StockOutService interface {
	GetStockOuts() ([]models.StockOut, error)
	CreateStockOut(stockOut *models.StockOut) (*models.StockOut, error)
	UpdateStockOutStatus(id uint, status models.StockOutStatus) (*models.StockOut, error)
}

type stockOutService struct {
	stockOutRepo repositories.StockOutRepository
	productRepo  repositories.ProductRepository
	stockLogRepo repositories.StockLogRepository
	db           *gorm.DB
}

func NewStockOutService(
	stockOutRepo repositories.StockOutRepository,
	productRepo repositories.ProductRepository,
	stockLogRepo repositories.StockLogRepository,
	db *gorm.DB,
) StockOutService {
	return &stockOutService{
		stockOutRepo: stockOutRepo,
		productRepo:  productRepo,
		stockLogRepo: stockLogRepo,
		db:           db,
	}
}

func (s *stockOutService) GetStockOuts() ([]models.StockOut, error) {
	return s.stockOutRepo.FindAll()
}

func (s *stockOutService) CreateStockOut(req *models.StockOut) (*models.StockOut, error) {
	var stockOut models.StockOut

	err := s.db.Transaction(func(tx *gorm.DB) error {
		txProductRepo := s.productRepo.WithTrx(tx)
		txStockOutRepo := s.stockOutRepo.WithTrx(tx)

		product, err := txProductRepo.FindByID(req.ProductID)
		if err != nil {
			return errors.New("product not found")
		}

		if product.AvailableStock < req.Quantity {
			return errors.New("insufficient available stock")
		}

		product.AvailableStock -= req.Quantity
		if err := txProductRepo.Update(product); err != nil {
			return err
		}

		stockOut = models.StockOut{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Status:    models.StockOutAllocated,
			Notes:     req.Notes,
		}
		if err := txStockOutRepo.Create(&stockOut); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &stockOut, nil
}

func (s *stockOutService) UpdateStockOutStatus(id uint, status models.StockOutStatus) (*models.StockOut, error) {
	var stockOut *models.StockOut
	var err error

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txStockOutRepo := s.stockOutRepo.WithTrx(tx)
		txProductRepo := s.productRepo.WithTrx(tx)
		txStockLogRepo := s.stockLogRepo.WithTrx(tx)

		stockOut, err = txStockOutRepo.FindByID(id)
		if err != nil {
			return err
		}

		if stockOut.Status == models.StockOutDone || stockOut.Status == models.StockOutCancelled {
			return errors.New("transaction is already final")
		}

		originalStatus := stockOut.Status
		stockOut.Status = status

		product, err := txProductRepo.FindByID(stockOut.ProductID)
		if err != nil {
			return err
		}

		switch status {
		case models.StockOutDone:
			if originalStatus != models.StockOutInProgress {
				return errors.New("can only complete transactions that are IN_PROGRESS")
			}
			previousStock := product.PhysicalStock
			product.PhysicalStock -= stockOut.Quantity
			if err := txProductRepo.Update(product); err != nil {
				return err
			}

			logEntry := models.StockLog{
				ProductID:       stockOut.ProductID,
				TransactionID:   stockOut.ID,
				TransactionType: models.TransactionStockOut,
				Quantity:        stockOut.Quantity,
				PreviousStock:   previousStock,
				CurrentStock:    product.PhysicalStock,
				Notes:           "Stock Out Completed",
			}
			if err := txStockLogRepo.Create(&logEntry); err != nil {
				return err
			}

		case models.StockOutCancelled:
			if originalStatus == models.StockOutAllocated || originalStatus == models.StockOutInProgress {
				product.AvailableStock += stockOut.Quantity
				if err := txProductRepo.Update(product); err != nil {
					return err
				}
			}

		case models.StockOutInProgress:
			if originalStatus != models.StockOutAllocated {
				return errors.New("can only start execution on ALLOCATED stock")
			}
		}

		if err := txStockOutRepo.Update(stockOut); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return stockOut, nil
}
