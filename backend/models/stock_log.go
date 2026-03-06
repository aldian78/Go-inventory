package models

import (
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionStockIn    TransactionType = "STOCK_IN"
	TransactionStockOut   TransactionType = "STOCK_OUT"
	TransactionAdjustment TransactionType = "ADJUSTMENT"
)

type StockLog struct {
	gorm.Model
	ProductID       uint            `json:"product_id"`
	Product         Product         `json:"product"`
	TransactionID   uint            `json:"transaction_id"` // ID of StockIn or StockOut
	TransactionType TransactionType `json:"transaction_type"`
	Quantity        int             `json:"quantity"`
	PreviousStock   int             `json:"previous_stock"`
	CurrentStock    int             `json:"current_stock"`
	Notes           string          `json:"notes"`
}
