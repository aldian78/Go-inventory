package repositories

import (
	"github.com/go-inventory/backend/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(name, sku, customer string) ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	WithTrx(trxHandle *gorm.DB) ProductRepository
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) WithTrx(trxHandle *gorm.DB) ProductRepository {
	if trxHandle == nil {
		return r
	}
	return &productRepository{db: trxHandle}
}

func (r *productRepository) FindAll(name, sku, customer string) ([]models.Product, error) {
	var products []models.Product
	db := r.db

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if sku != "" {
		db = db.Where("sku ILIKE ?", "%"+sku+"%")
	}
	if customer != "" {
		db = db.Where("customer ILIKE ?", "%"+customer+"%")
	}

	err := db.Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}
