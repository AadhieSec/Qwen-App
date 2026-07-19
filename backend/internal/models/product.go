package models

import (
"time"

"github.com/google/uuid"
"gorm.io/datatypes"
"gorm.io/gorm"
)

// Provider represents a shopping website provider
type Provider struct {
ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
Name        string         `gorm:"uniqueIndex;not null;size:100" json:"name"`
DisplayName string         `gorm:"size:255" json:"display_name"`
Domain      string         `gorm:"uniqueIndex;not null;size:255" json:"domain"`
BaseURL     string         `gorm:"size:512" json:"base_url"`
LogoURL     string         `gorm:"size:512" json:"logo_url"`
IsActive    bool           `gorm:"default:true" json:"is_active"`
Supports    datatypes.JSON `gorm:"type:jsonb" json:"supports"` // JSON array of supported features
CreatedAt   time.Time      `json:"created_at"`
UpdatedAt   time.Time      `json:"updated_at"`
}

func (Provider) TableName() string { return "providers" }

// Product represents a monitored product
type Product struct {
ID                uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
UserID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
ProviderID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"provider_id"`
URL               string         `gorm:"not null;size:2048;index" json:"url"`
Title             string         `gorm:"not null;size:1024" json:"title"`
Description       string         `gorm:"type:text" json:"description"`
Brand             string         `gorm:"size:255;index" json:"brand"`
Category          string         `gorm:"size:255;index" json:"category"`
ProductID         string         `gorm:"size:255;index" json:"product_id"` // Website's product ID
SKU               string         `gorm:"size:255" json:"sku"`
MainImage         string         `gorm:"size:2048" json:"main_image"`
Images            datatypes.JSON `gorm:"type:jsonb" json:"images"` // Array of image URLs
CurrentPrice      float64        `gorm:"type:decimal(10,2)" json:"current_price"`
Currency          string         `gorm:"size:10;default:'INR'" json:"currency"`
MRP               float64        `gorm:"type:decimal(10,2)" json:"mrp"`
Discount          float64        `gorm:"type:decimal(5,2)" json:"discount"` // Percentage
InStock           bool           `gorm:"default:true" json:"in_stock"`
StockStatus       string         `gorm:"size:50" json:"stock_status"` // in_stock, out_of_stock, low_stock
Rating            float64        `gorm:"type:decimal(3,2)" json:"rating"`
ReviewCount       int            `json:"review_count"`
Seller            string         `gorm:"size:255" json:"seller"`
Variants          datatypes.JSON `gorm:"type:jsonb" json:"variants"` // Sizes, colors, etc.
Coupons           datatypes.JSON `gorm:"type:jsonb" json:"coupons"`
DeliveryInfo      datatypes.JSON `gorm:"type:jsonb" json:"delivery_info"`
Metadata          datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
IsFavorite        bool           `gorm:"default:false" json:"is_favorite"`
IsMonitored       bool           `gorm:"default:true" json:"is_monitored"`
CheckInterval     int            `gorm:"default:3600" json:"check_interval"` // seconds
LastCheckedAt     *time.Time     `json:"last_checked_at"`
NextCheckAt       *time.Time     `json:"next_check_at"`
CreatedAt         time.Time      `json:"created_at"`
UpdatedAt         time.Time      `json:"updated_at"`
DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Product) TableName() string { return "products" }

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
Type        string    `gorm:"size:50;index" json:"type"` // size, color, storage, etc.
Value       string    `gorm:"size:100;index" json:"value"`
Available   bool      `gorm:"default:true" json:"available"`
Price       float64   `gorm:"type:decimal(10,2)" json:"price"`
SKU         string    `gorm:"size:255" json:"sku"`
CreatedAt   time.Time `json:"created_at"`
UpdatedAt   time.Time `json:"updated_at"`
}

func (ProductVariant) TableName() string { return "product_variants" }

// ProductImage represents additional product images
type ProductImage struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
URL       string    `gorm:"not null;size:2048" json:"url"`
Alt       string    `gorm:"size:512" json:"alt"`
Position  int       `gorm:"default:0" json:"position"`
CreatedAt time.Time `json:"created_at"`
}

func (ProductImage) TableName() string { return "product_images" }
