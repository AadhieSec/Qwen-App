package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"shopmonitor/internal/models"
)

// ProductRepository handles product data operations
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetByURL retrieves a product by URL for a specific user
func (r *ProductRepository) GetByURL(ctx context.Context, userID uuid.UUID, url string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("user_id = ? AND url = ?", userID, url).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetByUserID retrieves all products for a user
func (r *ProductRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	r.db.WithContext(ctx).Model(&models.Product{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// UpdatePrice updates product price and records history
func (r *ProductRepository) UpdatePrice(ctx context.Context, productID uuid.UUID, price, mrp, discount float64) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Update product
	if err := tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"current_price": price,
			"mrp":           mrp,
			"discount":      discount,
			"updated_at":    time.Now(),
		}).Error; err != nil {
		return err
	}

	// Record price history
	history := models.PriceHistory{
		ID:         uuid.New(),
		ProductID:  productID,
		Price:      price,
		Currency:   "INR",
		MRP:        mrp,
		Discount:   discount,
		RecordedAt: time.Now(),
		CreatedAt:  time.Now(),
	}
	if err := tx.Create(&history).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// UpdateStock updates product stock status
func (r *ProductRepository) UpdateStock(ctx context.Context, productID uuid.UUID, inStock bool, status string) error {
	tx := r.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"in_stock":        inStock,
			"stock_status":    status,
			"last_checked_at": time.Now(),
		}).Error; err != nil {
		return err
	}

	// Record stock history
	history := models.StockHistory{
		ID:          uuid.New(),
		ProductID:   productID,
		InStock:     inStock,
		StockStatus: status,
		RecordedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(&history).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// GetMonitoredProducts returns products that need monitoring
func (r *ProductRepository) GetMonitoredProducts(ctx context.Context, limit int) ([]*models.Product, error) {
	var products []*models.Product
	now := time.Now()

	err := r.db.WithContext(ctx).
		Where("is_monitored = ? AND is_paused = ? AND (next_check_at IS NULL OR next_check_at <= ?)", true, false, now).
		Limit(limit).
		Find(&products).Error

	return products, err
}

// UpdateCheckTimes updates last and next check timestamps
func (r *ProductRepository) UpdateCheckTimes(ctx context.Context, productID uuid.UUID, intervalSeconds int) error {
	now := time.Now()
	nextCheck := now.Add(time.Duration(intervalSeconds) * time.Second)

	return r.db.WithContext(ctx).Model(&models.Product{}).
		Where("id = ?", productID).
		Updates(map[string]interface{}{
			"last_checked_at": now,
			"next_check_at":   nextCheck,
		}).Error
}

// Search searches products by various criteria
func (r *ProductRepository) Search(ctx context.Context, userID uuid.UUID, query string, brand, category string, minPrice, maxPrice *float64, limit, offset int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	db := r.db.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID)

	if query != "" {
		db = db.Where("title ILIKE ? OR brand ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if brand != "" {
		db = db.Where("brand = ?", brand)
	}
	if category != "" {
		db = db.Where("category = ?", category)
	}
	if minPrice != nil {
		db = db.Where("current_price >= ?", *minPrice)
	}
	if maxPrice != nil {
		db = db.Where("current_price <= ?", *maxPrice)
	}

	db.Model(&models.Product{}).Count(&total)

	err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}
