package repositories

import (
	"time"

	"github.com/go-inventory/backend/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetStockInReport(startDate, endDate time.Time) ([]models.StockIn, error)
	GetStockOutReport(startDate, endDate time.Time) ([]models.StockOut, error)
	GetAdjustmentReport(startDate, endDate time.Time) ([]models.StockLog, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) applyDateFilter(db *gorm.DB, startDate, endDate time.Time) *gorm.DB {
	if !startDate.IsZero() {
		db = db.Where("created_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		endDate = endDate.Add(24*time.Hour - 1*time.Second) // End of the day
		db = db.Where("created_at <= ?", endDate)
	}
	return db
}

func (r *reportRepository) GetStockInReport(startDate, endDate time.Time) ([]models.StockIn, error) {
	var stockIns []models.StockIn
	db := r.db.Preload("Product").
		Where("status = ?", models.StockInDone)

	db = r.applyDateFilter(db, startDate, endDate)

	err := db.Order("created_at desc").Find(&stockIns).Error
	return stockIns, err
}

func (r *reportRepository) GetStockOutReport(startDate, endDate time.Time) ([]models.StockOut, error) {
	var stockOuts []models.StockOut
	db := r.db.Preload("Product").
		Where("status = ?", models.StockOutDone)

	db = r.applyDateFilter(db, startDate, endDate)

	err := db.Order("created_at desc").Find(&stockOuts).Error
	return stockOuts, err
}

func (r *reportRepository) GetAdjustmentReport(startDate, endDate time.Time) ([]models.StockLog, error) {
	var logs []models.StockLog
	db := r.db.Preload("Product").
		Where("transaction_type = ?", models.TransactionAdjustment)

	db = r.applyDateFilter(db, startDate, endDate)

	err := db.Order("created_at desc").Find(&logs).Error
	return logs, err
}
