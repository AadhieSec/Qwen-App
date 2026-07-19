package base

import (
"context"
"time"

"shopmonitor/internal/models"
)

// ProductData represents scraped product data
type ProductData struct {
Title       string                 `json:"title"`
Description string                 `json:"description"`
Brand       string                 `json:"brand"`
Category    string                 `json:"category"`
ProductID   string                 `json:"product_id"`
SKU         string                 `json:"sku"`
URL         string                 `json:"url"`
MainImage   string                 `json:"main_image"`
Images      []string               `json:"images"`
Price       float64                `json:"price"`
Currency    string                 `json:"currency"`
MRP         float64                `json:"mrp"`
Discount    float64                `json:"discount"`
InStock     bool                   `json:"in_stock"`
StockStatus string                 `json:"stock_status"`
Rating      float64                `json:"rating"`
ReviewCount int                    `json:"review_count"`
Seller      string                 `json:"seller"`
Variants    []Variant              `json:"variants"`
Coupons     []Coupon               `json:"coupons"`
Delivery    []DeliveryInfo         `json:"delivery"`
Metadata    map[string]interface{} `json:"metadata"`
}

// Variant represents a product variant (size, color, etc.)
type Variant struct {
Type      string  `json:"type"`
Value     string  `json:"value"`
Available bool    `json:"available"`
Price     float64 `json:"price"`
SKU       string  `json:"sku"`
}

// Coupon represents a coupon or offer
type Coupon struct {
Code        string     `json:"code"`
Description string     `json:"description"`
Discount    string     `json:"discount"`
ValidUntil  *time.Time `json:"valid_until"`
IsActive    bool       `json:"is_active"`
}

// DeliveryInfo represents delivery information for a pincode
type DeliveryInfo struct {
Pincode       string  `json:"pincode"`
Available     bool    `json:"available"`
DeliveryType  string  `json:"delivery_type"`
EstimatedDays int     `json:"estimated_days"`
Fee           float64 `json:"fee"`
}

// Provider defines the interface that all shopping providers must implement
type Provider interface {
// GetInfo returns provider metadata
GetInfo() ProviderInfo

// CanHandle checks if this provider can handle the given URL
CanHandle(url string) bool

// Discover extracts product information from a URL without full scraping
Discover(ctx context.Context, url string) (*ProductData, error)

// Monitor performs a full product monitoring check
Monitor(ctx context.Context, product *models.Product) (*ProductData, error)

// FetchPrice gets the current price for a product
FetchPrice(ctx context.Context, url string) (float64, error)

// FetchVariants gets available variants (sizes, colors, etc.)
FetchVariants(ctx context.Context, url string) ([]Variant, error)

// FetchCoupons gets available coupons for the product
FetchCoupons(ctx context.Context, url string) ([]Coupon, error)

// FetchDelivery checks delivery availability for a pincode
FetchDelivery(ctx context.Context, url, pincode string) (*DeliveryInfo, error)

// FetchStock checks stock availability
FetchStock(ctx context.Context, url string) (bool, string, error)

// FetchMetadata gets additional product metadata
FetchMetadata(ctx context.Context, url string) (map[string]interface{}, error)

// HealthCheck verifies the provider is working
HealthCheck(ctx context.Context) error

// Close cleans up any resources
Close() error
}

// ProviderInfo contains metadata about a provider
type ProviderInfo struct {
Name        string   `json:"name"`
DisplayName string   `json:"display_name"`
Domain      string   `json:"domain"`
BaseURL     string   `json:"base_url"`
LogoURL     string   `json:"logo_url"`
Supports    []string `json:"supports"` // List of supported features
}

// MonitoringMethod represents different approaches to fetch data
type MonitoringMethod int

const (
MethodAPI MonitoringMethod = iota
MethodGraphQL
MethodNetworkAPI
MethodHTMLParsing
MethodBrowserAutomation
)

// StrategyResult holds the result of method selection
type StrategyResult struct {
Method MonitoringMethod
URL    string
Headers map[string]string
Data   interface{}
}
