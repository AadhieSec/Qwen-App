package amazon

import (
"context"
"encoding/json"
"fmt"
"regexp"
"strings"
"time"

"shopmonitor/internal/models"
"shopmonitor/internal/provider/base"
)

// AmazonProvider implements the Provider interface for Amazon
type AmazonProvider struct {
*base.BaseProvider
info base.ProviderInfo
}

// NewAmazonProvider creates a new Amazon provider
func NewAmazonProvider(userAgent string, timeout time.Duration) *AmazonProvider {
return &AmazonProvider{
BaseProvider: base.NewBaseProvider(userAgent, timeout),
info: base.ProviderInfo{
Name:        "amazon",
DisplayName: "Amazon",
Domain:      "amazon.in",
BaseURL:     "https://www.amazon.in",
LogoURL:     "https://logo.clearbit.com/amazon.in",
Supports:    []string{"price", "stock", "variants", "coupons", "delivery"},
},
}
}

// GetInfo returns provider metadata
func (p *AmazonProvider) GetInfo() base.ProviderInfo {
return p.info
}

// CanHandle checks if this provider can handle the given URL
func (p *AmazonProvider) CanHandle(url string) bool {
domains := []string{
"amazon.in",
"amazon.com",
"amazon.co.uk",
"amazon.de",
"amazon.fr",
"amazon.co.jp",
"amazon.ca",
"amazon.com.au",
}

for _, domain := range domains {
if strings.Contains(url, domain) {
return true
}
}
return false
}

// Discover extracts product information from a URL
func (p *AmazonProvider) Discover(ctx context.Context, url string) (*base.ProductData, error) {
productID, err := p.extractProductID(url)
if err != nil {
return nil, err
}

return &base.ProductData{
ProductID: productID,
URL:       url,
Currency:  "INR",
Metadata: map[string]interface{}{
"provider": "amazon",
},
}, nil
}

// Monitor performs a full product monitoring check
func (p *AmazonProvider) Monitor(ctx context.Context, product *models.Product) (*base.ProductData, error) {
html, err := p.FetchURL(ctx, product.URL)
if err != nil {
return nil, fmt.Errorf("failed to fetch URL: %w", err)
}

data := &base.ProductData{
URL:      product.URL,
Currency: "INR",
Metadata: map[string]interface{}{
"provider": "amazon",
},
}

// Extract price using regex patterns
pricePatterns := []string{
`₹\s*([\d,]+\.?\d*)`,
`class="a-price-whole">([\d,]+)`,
`"price":"([\d.]+)"`,
}

for _, pattern := range pricePatterns {
re := regexp.MustCompile(pattern)
matches := re.FindStringSubmatch(html)
if len(matches) > 1 {
priceStr := strings.ReplaceAll(matches[1], ",", "")
fmt.Sscanf(priceStr, "%f", &data.Price)
break
}
}

// Extract MRP
mrpPattern := regexp.MustCompile(`(?:MRP|List price):.*?₹\s*([\d,]+\.?\d*)`)
if matches := mrpPattern.FindStringSubmatch(html); len(matches) > 1 {
mrpStr := strings.ReplaceAll(matches[1], ",", "")
fmt.Sscanf(mrpStr, "%f", &data.MRP)
}

// Calculate discount
if data.Price > 0 && data.MRP > 0 {
data.Discount = ((data.MRP - data.Price) / data.MRP) * 100
}

// Check stock
data.InStock = !strings.Contains(html, "Currently unavailable") &&
!strings.Contains(html, "Out of stock")
data.StockStatus = p.DetermineStockStatus(data.InStock, "")

// Extract title
titlePattern := regexp.MustCompile(`<title>([^<]+)</title>`)
if matches := titlePattern.FindStringSubmatch(html); len(matches) > 1 {
data.Title = strings.TrimSpace(matches[1])
}

// Extract product ID
if productID, err := p.extractProductID(product.URL); err == nil {
data.ProductID = productID
}

return data, nil
}

// FetchPrice gets the current price for a product
func (p *AmazonProvider) FetchPrice(ctx context.Context, url string) (float64, error) {
html, err := p.FetchURL(ctx, url)
if err != nil {
return 0, err
}

pricePatterns := []string{
`₹\s*([\d,]+\.?\d*)`,
`class="a-price-whole">([\d,]+)`,
`"price":"([\d.]+)"`,
}

for _, pattern := range pricePatterns {
re := regexp.MustCompile(pattern)
matches := re.FindStringSubmatch(html)
if len(matches) > 1 {
priceStr := strings.ReplaceAll(matches[1], ",", "")
var price float64
fmt.Sscanf(priceStr, "%f", &price)
return price, nil
}
}

return 0, fmt.Errorf("price not found")
}

// FetchVariants gets available variants
func (p *AmazonProvider) FetchVariants(ctx context.Context, url string) ([]base.Variant, error) {
html, err := p.FetchURL(ctx, url)
if err != nil {
return nil, err
}

var variants []base.Variant

// Look for size variants
sizePattern := regexp.MustCompile(`data-asin="([A-Z0-9]+)"[^>]*>([^<]+)<`)
matches := sizePattern.FindAllStringSubmatch(html, -1)
for _, match := range matches {
if len(match) > 2 {
variants = append(variants, base.Variant{
Type:      "size",
Value:     strings.TrimSpace(match[2]),
Available: true,
SKU:       match[1],
})
}
}

return variants, nil
}

// FetchCoupons gets available coupons
func (p *AmazonProvider) FetchCoupons(ctx context.Context, url string) ([]base.Coupon, error) {
html, err := p.FetchURL(ctx, url)
if err != nil {
return nil, err
}

var coupons []base.Coupon

// Look for coupon clues
couponPattern := regexp.MustCompile(`Save.*?(?:₹|Rs\.?\s*)(\d+%?)`)
matches := couponPattern.FindAllStringSubmatch(html, -1)
for _, match := range matches {
if len(match) > 1 {
coupons = append(coupons, base.Coupon{
Code:        "AUTO",
Description: fmt.Sprintf("Save %s", match[1]),
Discount:    match[1],
IsActive:    true,
})
}
}

return coupons, nil
}

// FetchDelivery checks delivery availability
func (p *AmazonProvider) FetchDelivery(ctx context.Context, url, pincode string) (*base.DeliveryInfo, error) {
// For now, return basic info - would need actual pincode check API
return &base.DeliveryInfo{
Pincode:       pincode,
Available:     true,
DeliveryType:  "standard",
EstimatedDays: 3,
Fee:           0,
}, nil
}

// FetchStock checks stock availability
func (p *AmazonProvider) FetchStock(ctx context.Context, url string) (bool, string, error) {
html, err := p.FetchURL(ctx, url)
if err != nil {
return false, "", err
}

inStock := !strings.Contains(html, "Currently unavailable") &&
!strings.Contains(html, "Out of stock")

status := "in_stock"
if !inStock {
status = "out_of_stock"
}

return inStock, status, nil
}

// FetchMetadata gets additional product metadata
func (p *AmazonProvider) FetchMetadata(ctx context.Context, url string) (map[string]interface{}, error) {
html, err := p.FetchURL(ctx, url)
if err != nil {
return nil, err
}

metadata := make(map[string]interface{})

// Extract ASIN
if asin, err := p.extractProductID(url); err == nil {
metadata["asin"] = asin
}

// Extract brand
brandPattern := regexp.MustCompile(`Brand:\s*<span[^>]*>([^<]+)`)
if matches := brandPattern.FindStringSubmatch(html); len(matches) > 1 {
metadata["brand"] = strings.TrimSpace(matches[1])
}

return metadata, nil
}

// HealthCheck verifies the provider is working
func (p *AmazonProvider) HealthCheck(ctx context.Context) error {
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

_, err := p.FetchURL(ctx, "https://www.amazon.in")
return err
}

// Close cleans up resources
func (p *AmazonProvider) Close() error {
return nil
}

// extractProductID extracts Amazon product ID (ASIN) from URL
func (p *AmazonProvider) extractProductID(url string) (string, error) {
patterns := []string{
`/dp/([A-Z0-9]{10})`,
`/gp/product/([A-Z0-9]{10})`,
`asin=([A-Z0-9]{10})`,
}

for _, pattern := range patterns {
re := regexp.MustCompile(pattern)
matches := re.FindStringSubmatch(url)
if len(matches) > 1 {
return matches[1], nil
}
}

return "", fmt.Errorf("could not extract product ID from URL")
}

// Ensure AmazonProvider implements base.Provider
var _ base.Provider = (*AmazonProvider)(nil)
