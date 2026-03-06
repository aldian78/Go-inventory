package services

import (
	"errors"
	"time"

	"github.com/go-inventory/backend/repositories"
)

type ReportService interface {
	GetTransactionReport(transType string, startDate, endDate time.Time) (interface{}, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{reportRepo}
}

func (s *reportService) GetTransactionReport(transType string, startDate, endDate time.Time) (interface{}, error) {
	switch transType {
	case "stock_in":
		return s.reportRepo.GetStockInReport(startDate, endDate)
	case "stock_out":
		return s.reportRepo.GetStockOutReport(startDate, endDate)
	case "adjustment":
		return s.reportRepo.GetAdjustmentReport(startDate, endDate)
	default:
		return nil, errors.New("invalid report type. Use 'stock_in', 'stock_out', or 'adjustment'")
	}
}
