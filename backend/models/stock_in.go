package models

import (
	"gorm.io/gorm"
)

type StockInStatus string

const (
	StockInCreated    StockInStatus = "CREATED"
	StockInInProgress StockInStatus = "IN_PROGRESS"
	StockInDone       StockInStatus = "DONE"
	StockInCancelled  StockInStatus = "CANCELLED"
)

type StockIn struct {
	gorm.Model
	ProductID uint          `json:"product_id"`
	Product   Product       `json:"product"`
	Quantity  int           `json:"quantity"`
	Status    StockInStatus `json:"status" gorm:"default:'CREATED'"`
	Notes     string        `json:"notes"`
}
