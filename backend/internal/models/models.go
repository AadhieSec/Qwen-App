package models

import (
"time"

"github.com/google/uuid"
)

// Tag represents a user-defined tag for products
type Tag struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
Name      string    `gorm:"size:100;index" json:"name"`
Color     string    `gorm:"size:20" json:"color"`
CreatedAt time.Time `json:"created_at"`
}

func (Tag) TableName() string { return "tags" }

// ProductTag links products to tags
type ProductTag struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"product_id"`
TagID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tag_id"`
CreatedAt time.Time `json:"created_at"`
}

func (ProductTag) TableName() string { return "product_tags" }

// Wishlist represents user's wishlist items
type Wishlist struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
Name      string    `gorm:"size:255" json:"name"`
IsPublic  bool      `gorm:"default:false" json:"is_public"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

func (Wishlist) TableName() string { return "wishlists" }

// WishlistItem represents an item in a wishlist
type WishlistItem struct {
ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
WishlistID uuid.UUID `gorm:"type:uuid;not null;index" json:"wishlist_id"`
ProductID  uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
AddedAt    time.Time `json:"added_at"`
}

func (WishlistItem) TableName() string { return "wishlist_items" }

// UserSetting represents user preferences
type UserSetting struct {
ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
UserID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
Key       string         `gorm:"size:100;index" json:"key"`
Value     interface{}    `gorm:"type:jsonb" json:"value"`
UpdatedAt time.Time      `json:"updated_at"`
}

func (UserSetting) TableName() string { return "user_settings" }

// SavedSearch represents a saved product search
type SavedSearch struct {
ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
UserID      uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
Name        string    `gorm:"size:255" json:"name"`
Query       string    `gorm:"type:text" json:"query"`
Filters     interface{} `gorm:"type:jsonb" json:"filters"`
CheckInterval int     `gorm:"default:3600" json:"check_interval"`
IsActive    bool      `gorm:"default:true" json:"is_active"`
CreatedAt   time.Time `json:"created_at"`
UpdatedAt   time.Time `json:"updated_at"`
}

func (SavedSearch) TableName() string { return "saved_searches" }

// AuditLog tracks important system events
type AuditLog struct {
ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
UserID     *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
Action     string    `gorm:"size:100;index" json:"action"`
Resource   string    `gorm:"size:100" json:"resource"`
ResourceID *uuid.UUID `gorm:"type:uuid" json:"resource_id,omitempty"`
Details    interface{} `gorm:"type:jsonb" json:"details"`
IPAddress  string    `gorm:"size:45" json:"ip_address"`
UserAgent  string    `gorm:"size:512" json:"user_agent"`
CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }
