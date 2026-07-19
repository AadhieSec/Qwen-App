package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationChannel represents user's notification preferences
type NotificationChannel struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Type       string         `gorm:"size:50;index" json:"type"`  // desktop, telegram, discord, slack, email, webhook, push
	Identifier string         `gorm:"size:512" json:"identifier"` // bot token, webhook URL, email, etc.
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	Settings   datatypes.JSON `gorm:"type:jsonb" json:"settings"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (NotificationChannel) TableName() string { return "notification_channels" }

// NotificationPreference represents what events user wants to be notified about
type NotificationPreference struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	PriceDrop         bool      `gorm:"default:true" json:"price_drop"`
	PriceTarget       bool      `gorm:"default:true" json:"price_target"`
	DiscountIncrease  bool      `gorm:"default:true" json:"discount_increase"`
	BackInStock       bool      `gorm:"default:true" json:"back_in_stock"`
	OutOfStock        bool      `gorm:"default:false" json:"out_of_stock"`
	SizeAvailable     bool      `gorm:"default:true" json:"size_available"`
	ColorAvailable    bool      `gorm:"default:true" json:"color_available"`
	CouponAvailable   bool      `gorm:"default:true" json:"coupon_available"`
	DeliveryAvailable bool      `gorm:"default:true" json:"delivery_available"`
	FlashSale         bool      `gorm:"default:true" json:"flash_sale"`
	PageChange        bool      `gorm:"default:false" json:"page_change"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (NotificationPreference) TableName() string { return "notification_preferences" }

// Notification represents a generated notification
type Notification struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID    *uuid.UUID     `gorm:"type:uuid;index" json:"product_id,omitempty"`
	Type         string         `gorm:"size:50;index" json:"type"` // price_drop, back_in_stock, size_available, etc.
	Title        string         `gorm:"size:255" json:"title"`
	Message      string         `gorm:"type:text" json:"message"`
	Data         datatypes.JSON `gorm:"type:jsonb" json:"data"` // Additional context
	ImageURL     string         `gorm:"size:2048" json:"image_url"`
	ActionURL    string         `gorm:"size:2048" json:"action_url"`
	IsRead       bool           `gorm:"default:false;index" json:"is_read"`
	Priority     string         `gorm:"size:20;default:normal" json:"priority"` // low, normal, high, critical
	SentChannels datatypes.JSON `gorm:"type:jsonb" json:"sent_channels"`        // List of channels where sent
	SentAt       *time.Time     `json:"sent_at"`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string { return "notifications" }

// NotificationLog tracks notification delivery attempts
type NotificationLog struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	NotificationID uuid.UUID `gorm:"type:uuid;not null;index" json:"notification_id"`
	ChannelType    string    `gorm:"size:50" json:"channel_type"`
	Status         string    `gorm:"size:20;index" json:"status"` // pending, sent, failed, skipped
	ErrorMessage   string    `gorm:"type:text" json:"error_message,omitempty"`
	ResponseData   string    `gorm:"type:text" json:"response_data,omitempty"`
	AttemptedAt    time.Time `gorm:"index" json:"attempted_at"`
	CreatedAt      time.Time `json:"created_at"`
}

func (NotificationLog) TableName() string { return "notification_logs" }
