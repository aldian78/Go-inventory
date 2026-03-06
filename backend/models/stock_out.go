package models

import (
	"gorm.io/gorm"
)

type StockOutStatus string

const (
	StockOutDraft      StockOutStatus = "DRAFT"
	StockOutAllocated  StockOutStatus = "ALLOCATED"
	StockOutInProgress StockOutStatus = "IN_PROGRESS"
	StockOutDone       StockOutStatus = "DONE"
	StockOutCancelled  StockOutStatus = "CANCELLED"
)

type StockOut struct {
	gorm.Model
	ProductID uint           `json:"product_id"`
	Product   Product        `json:"product"`
	Quantity  int            `json:"quantity"`
	Status    StockOutStatus `json:"status" gorm:"default:'DRAFT'"`
	Notes     string         `json:"notes"`
}
