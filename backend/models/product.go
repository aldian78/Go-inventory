package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name           string `json:"name"`
	SKU            string `json:"sku" gorm:"unique"`
	Customer       string `json:"customer"` // Added Customer field
	PhysicalStock  int    `json:"physical_stock"`
	AvailableStock int    `json:"available_stock"`
}
