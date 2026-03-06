package repositories

import (
	"github.com/go-inventory/backend/models"
	"gorm.io/gorm"
)

type StockInRepository interface {
	FindAll() ([]models.StockIn, error)
	FindByID(id uint) (*models.StockIn, error)
	Create(stockIn *models.StockIn) error
	Update(stockIn *models.StockIn) error
	WithTrx(trxHandle *gorm.DB) StockInRepository
}

type stockInRepository struct {
	db *gorm.DB
}

func NewStockInRepository(db *gorm.DB) StockInRepository {
	return &stockInRepository{db}
}

func (r *stockInRepository) WithTrx(trxHandle *gorm.DB) StockInRepository {
	if trxHandle == nil {
		return r
	}
	return &stockInRepository{db: trxHandle}
}

func (r *stockInRepository) FindAll() ([]models.StockIn, error) {
	var stockIns []models.StockIn
	err := r.db.Preload("Product").Find(&stockIns).Error
	return stockIns, err
}

func (r *stockInRepository) FindByID(id uint) (*models.StockIn, error) {
	var stockIn models.StockIn
	err := r.db.First(&stockIn, id).Error
	return &stockIn, err
}

func (r *stockInRepository) Create(stockIn *models.StockIn) error {
	return r.db.Create(stockIn).Error
}

func (r *stockInRepository) Update(stockIn *models.StockIn) error {
	return r.db.Save(stockIn).Error
}
