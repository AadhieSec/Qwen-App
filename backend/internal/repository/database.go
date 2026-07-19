package repository

import (
	"fmt"

	"shopmonitor/internal/config"
	"shopmonitor/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the GORM database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &Database{DB: db}, nil
}

// AutoMigrate runs auto migration for all models
func (d *Database) AutoMigrate() error {
	models := []interface{}{
		&models.User{},
		&models.Provider{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},
		&models.PriceHistory{},
		&models.StockHistory{},
		&models.VariantHistory{},
		&models.CouponHistory{},
		&models.DeliveryHistory{},
		&models.NotificationChannel{},
		&models.NotificationPreference{},
		&models.Notification{},
		&models.NotificationLog{},
		&models.Job{},
		&models.Worker{},
		&models.MonitorConfig{},
		&models.Tag{},
		&models.ProductTag{},
		&models.Wishlist{},
		&models.WishlistItem{},
		&models.UserSetting{},
		&models.SavedSearch{},
		&models.AuditLog{},
	}

	for _, model := range models {
		if err := d.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("migration failed for %T: %w", model, err)
		}
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks database connectivity
func (d *Database) HealthCheck() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
