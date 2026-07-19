package repository

import (
"context"
"errors"
"time"

"github.com/google/uuid"
"shopmonitor/internal/models"
"gorm.io/gorm"
)

// UserRepository handles user data operations
type UserRepository struct {
db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
var user models.User
err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, nil
}
return nil, err
}
return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
var user models.User
err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, nil
}
return nil, err
}
return &user, nil
}

// GetByProviderID retrieves a user by provider ID
func (r *UserRepository) GetByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
var user models.User
err := r.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, nil
}
return nil, err
}
return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
return r.db.WithContext(ctx).Save(user).Error
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
now := time.Now()
return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

// UpdateRefreshToken updates the refresh token
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, token string) error {
return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("refresh_token", token).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// Exists checks if a user exists by email
func (r *UserRepository) Exists(ctx context.Context, email string) (bool, error) {
var count int64
err := r.db.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
return count > 0, err
}
