package models

import (
"time"

"github.com/google/uuid"
"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
Email           string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
Password        string         `gorm:"not null;size:255" json:"-"`
Name            string         `gorm:"size:255" json:"name"`
Avatar          string         `gorm:"size:512" json:"avatar"`
EmailVerified   bool           `gorm:"default:false" json:"email_verified"`
Provider        string         `gorm:"size:50" json:"provider"` // local, google, github
ProviderID      string         `gorm:"size:255" json:"provider_id"`
RefreshToken    string         `gorm:"size:512" json:"-"`
LastLoginAt     *time.Time     `json:"last_login_at"`
CreatedAt       time.Time      `json:"created_at"`
UpdatedAt       time.Time      `json:"updated_at"`
DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate generates UUID before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) error {
if u.ID == uuid.Nil {
u.ID = uuid.New()
}
return nil
}

// TableName returns the table name for User model
func (User) TableName() string {
return "users"
}
