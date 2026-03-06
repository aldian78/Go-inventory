package repositories

import (
	"github.com/go-inventory/backend/models"
	"gorm.io/gorm"
)

type StockOutRepository interface {
	FindAll() ([]models.StockOut, error)
	FindByID(id uint) (*models.StockOut, error)
	Create(stockOut *models.StockOut) error
	Update(stockOut *models.StockOut) error
	WithTrx(trxHandle *gorm.DB) StockOutRepository
}

type stockOutRepository struct {
	db *gorm.DB
}

func NewStockOutRepository(db *gorm.DB) StockOutRepository {
	return &stockOutRepository{db}
}

func (r *stockOutRepository) WithTrx(trxHandle *gorm.DB) StockOutRepository {
	if trxHandle == nil {
		return r
	}
	return &stockOutRepository{db: trxHandle}
}

func (r *stockOutRepository) FindAll() ([]models.StockOut, error) {
	var stockOuts []models.StockOut
	err := r.db.Preload("Product").Find(&stockOuts).Error
	return stockOuts, err
}

func (r *stockOutRepository) FindByID(id uint) (*models.StockOut, error) {
	var stockOut models.StockOut
	err := r.db.First(&stockOut, id).Error
	return &stockOut, err
}

func (r *stockOutRepository) Create(stockOut *models.StockOut) error {
	return r.db.Create(stockOut).Error
}

func (r *stockOutRepository) Update(stockOut *models.StockOut) error {
	return r.db.Save(stockOut).Error
}
