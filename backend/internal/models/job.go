package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Job represents a monitoring job in the queue
type Job struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Type         string         `gorm:"size:100;not null;index" json:"type"`         // product_check, notification_send, etc.
	Priority     int            `gorm:"default:0;index" json:"priority"`             // Higher = more urgent
	Status       string         `gorm:"size:20;default:pending;index" json:"status"` // pending, processing, completed, failed, dead_letter
	Payload      datatypes.JSON `gorm:"type:jsonb" json:"payload"`
	RetryCount   int            `gorm:"default:0" json:"retry_count"`
	MaxRetries   int            `gorm:"default:3" json:"max_retries"`
	ScheduledAt  time.Time      `gorm:"index" json:"scheduled_at"`
	StartedAt    *time.Time     `json:"started_at"`
	CompletedAt  *time.Time     `json:"completed_at"`
	ErrorMessage string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (Job) TableName() string { return "jobs" }

// Worker represents a worker instance
type Worker struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name          string     `gorm:"size:255;uniqueIndex" json:"name"`
	Type          string     `gorm:"size:100" json:"type"`
	Status        string     `gorm:"size:20;default:idle" json:"status"` // idle, busy, offline
	CurrentJobID  *uuid.UUID `gorm:"type:uuid" json:"current_job_id,omitempty"`
	JobsProcessed int64      `json:"jobs_processed"`
	JobsFailed    int64      `json:"jobs_failed"`
	LastHeartbeat time.Time  `gorm:"index" json:"last_heartbeat"`
	StartedAt     time.Time  `json:"started_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Worker) TableName() string { return "workers" }

// MonitorConfig represents user's monitoring configuration for a product
type MonitorConfig struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	ProductID        uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"product_id"`
	TargetPrice      *float64       `gorm:"type:decimal(10,2)" json:"target_price"`
	TargetDiscount   *float64       `gorm:"type:decimal(5,2)" json:"target_discount"`
	MaxPrice         *float64       `gorm:"type:decimal(10,2)" json:"max_price"`
	MinDiscount      *float64       `gorm:"type:decimal(5,2)" json:"min_discount"`
	DesiredSizes     datatypes.JSON `gorm:"type:jsonb" json:"desired_sizes"`  // Array of sizes
	DesiredColors    datatypes.JSON `gorm:"type:jsonb" json:"desired_colors"` // Array of colors
	DesiredVariants  datatypes.JSON `gorm:"type:jsonb" json:"desired_variants"`
	SellerPreference string         `gorm:"size:255" json:"seller_preference"`
	DeliveryPincode  string         `gorm:"size:10" json:"delivery_pincode"`
	CheckInterval    int            `gorm:"default:3600" json:"check_interval"` // seconds
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	IsPaused         bool           `gorm:"default:false" json:"is_paused"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (MonitorConfig) TableName() string { return "monitor_configs" }
