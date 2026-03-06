package repositories

import (
	"github.com/go-inventory/backend/models"
	"gorm.io/gorm"
)

type StockLogRepository interface {
	Create(log *models.StockLog) error
	WithTrx(trxHandle *gorm.DB) StockLogRepository
}

type stockLogRepository struct {
	db *gorm.DB
}

func NewStockLogRepository(db *gorm.DB) StockLogRepository {
	return &stockLogRepository{db}
}

func (r *stockLogRepository) WithTrx(trxHandle *gorm.DB) StockLogRepository {
	if trxHandle == nil {
		return r
	}
	return &stockLogRepository{db: trxHandle}
}

func (r *stockLogRepository) Create(log *models.StockLog) error {
	return r.db.Create(log).Error
}
