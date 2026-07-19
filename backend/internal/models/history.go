package models

import (
"time"

"github.com/google/uuid"
)

// PriceHistory tracks price changes over time
type PriceHistory struct {
ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
Currency    string    `gorm:"size:10" json:"currency"`
MRP         float64   `gorm:"type:decimal(10,2)" json:"mrp"`
Discount    float64   `gorm:"type:decimal(5,2)" json:"discount"`
InStock     bool      `gorm:"default:true" json:"in_stock"`
Seller      string    `gorm:"size:255" json:"seller"`
RecordedAt  time.Time `gorm:"index" json:"recorded_at"`
CreatedAt   time.Time `json:"created_at"`
}

func (PriceHistory) TableName() string { return "price_history" }

// StockHistory tracks stock availability changes
type StockHistory struct {
ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
InStock     bool      `gorm:"default:true" json:"in_stock"`
StockStatus string    `gorm:"size:50" json:"stock_status"`
Quantity    *int      `json:"quantity,omitempty"` // If available
RecordedAt  time.Time `gorm:"index" json:"recorded_at"`
CreatedAt   time.Time `json:"created_at"`
}

func (StockHistory) TableName() string { return "stock_history" }

// VariantHistory tracks variant availability changes
type VariantHistory struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
VariantID uuid.UUID `gorm:"type:uuid;index" json:"variant_id"`
Type      string    `gorm:"size:50" json:"type"`
Value     string    `gorm:"size:100" json:"value"`
Available bool      `gorm:"default:true" json:"available"`
RecordedAt time.Time `gorm:"index" json:"recorded_at"`
CreatedAt  time.Time `json:"created_at"`
}

func (VariantHistory) TableName() string { return "variant_history" }

// CouponHistory tracks coupon availability
type CouponHistory struct {
ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID   uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
Code        string    `gorm:"size:100;index" json:"code"`
Description string    `gorm:"size:512" json:"description"`
Discount    string    `gorm:"size:100" json:"discount"`
IsActive    bool      `gorm:"default:true" json:"is_active"`
ValidUntil  *time.Time `json:"valid_until"`
RecordedAt  time.Time `gorm:"index" json:"recorded_at"`
CreatedAt   time.Time `json:"created_at"`
}

func (CouponHistory) TableName() string { return "coupon_history" }

// DeliveryHistory tracks delivery availability for pincodes
type DeliveryHistory struct {
ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID    uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
Pincode      string    `gorm:"size:10;index" json:"pincode"`
Available    bool      `gorm:"default:false" json:"available"`
DeliveryType string    `gorm:"size:50" json:"delivery_type"` // standard, express, same_day
EstimatedDays int      `json:"estimated_days"`
Fee          float64   `gorm:"type:decimal(10,2)" json:"fee"`
RecordedAt   time.Time `gorm:"index" json:"recorded_at"`
CreatedAt    time.Time `json:"created_at"`
}

func (DeliveryHistory) TableName() string { return "delivery_history" }
