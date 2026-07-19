package flipkart

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"shopmonitor/internal/models"
	"shopmonitor/internal/provider/base"
)

// FlipkartProvider implements the Provider interface for Flipkart
type FlipkartProvider struct {
	*base.BaseProvider
	info base.ProviderInfo
}

// NewFlipkartProvider creates a new Flipkart provider
func NewFlipkartProvider(userAgent string, timeout time.Duration) *FlipkartProvider {
	return &FlipkartProvider{
		BaseProvider: base.NewBaseProvider(userAgent, timeout),
		info: base.ProviderInfo{
			Name:        "flipkart",
			DisplayName: "Flipkart",
			Domain:      "flipkart.com",
			BaseURL:     "https://www.flipkart.com",
			LogoURL:     "https://logo.clearbit.com/flipkart.com",
			Supports:    []string{"price", "stock", "variants", "coupons", "delivery"},
		},
	}
}

// GetInfo returns provider metadata
func (p *FlipkartProvider) GetInfo() base.ProviderInfo {
	return p.info
}

// CanHandle checks if this provider can handle the given URL
func (p *FlipkartProvider) CanHandle(url string) bool {
	domains := []string{
		"flipkart.com",
		"fkrt.it",
	}

	for _, domain := range domains {
		if strings.Contains(url, domain) {
			return true
		}
	}
	return false
}

// Discover extracts product information from a URL
func (p *FlipkartProvider) Discover(ctx context.Context, url string) (*base.ProductData, error) {
	productID, err := p.extractProductID(url)
	if err != nil {
		return nil, err
	}

	return &base.ProductData{
		ProductID: productID,
		URL:       url,
		Currency:  "INR",
		Metadata: map[string]interface{}{
			"provider": "flipkart",
		},
	}, nil
}

// Monitor performs a full product monitoring check
func (p *FlipkartProvider) Monitor(ctx context.Context, product *models.Product) (*base.ProductData, error) {
	html, err := p.FetchURL(ctx, product.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}

	data := &base.ProductData{
		URL:      product.URL,
		Currency: "INR",
		Metadata: map[string]interface{}{
			"provider": "flipkart",
		},
	}

	// Extract price using regex patterns
	pricePatterns := []string{
		`₹\s*([\d,]+\.?\d*)`,
		`class="_30jeq3".*?>([\d,]+\.?\d*)`,
		`"price":\s*"([\d.]+)"`,
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
	mrpPattern := regexp.MustCompile(`(?:MRP|List Price).*?₹\s*([\d,]+\.?\d*)`)
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
		!strings.Contains(html, "Out of stock") &&
		!strings.Contains(html, "Sold Out")
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
func (p *FlipkartProvider) FetchPrice(ctx context.Context, url string) (float64, error) {
	html, err := p.FetchURL(ctx, url)
	if err != nil {
		return 0, err
	}

	pricePatterns := []string{
		`₹\s*([\d,]+\.?\d*)`,
		`class="_30jeq3".*?>([\d,]+\.?\d*)`,
		`"price":\s*"([\d.]+)"`,
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
func (p *FlipkartProvider) FetchVariants(ctx context.Context, url string) ([]base.Variant, error) {
	html, err := p.FetchURL(ctx, url)
	if err != nil {
		return nil, err
	}

	var variants []base.Variant

	// Look for size/color variants in selection elements
	variantPattern := regexp.MustCompile(`(?:value|data-value)=["']([^"']+)["'][^>]*(?:selected|active)`)
	matches := variantPattern.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) > 1 {
			value := strings.TrimSpace(match[1])
			if value != "" {
				variants = append(variants, base.Variant{
					Type:      "variant",
					Value:     value,
					Available: true,
				})
			}
		}
	}

	return variants, nil
}

// FetchCoupons gets available coupons
func (p *FlipkartProvider) FetchCoupons(ctx context.Context, url string) ([]base.Coupon, error) {
	html, err := p.FetchURL(ctx, url)
	if err != nil {
		return nil, err
	}

	var coupons []base.Coupon

	// Look for coupon clues
	couponPattern := regexp.MustCompile(`(?:Bank Offer|Coupon|Discount).*?(?:₹|Rs\.?\s*)(\d+%?)`)
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
func (p *FlipkartProvider) FetchDelivery(ctx context.Context, url, pincode string) (*base.DeliveryInfo, error) {
	// For now, return basic info - would need actual pincode check API
	return &base.DeliveryInfo{
		Pincode:       pincode,
		Available:     true,
		DeliveryType:  "standard",
		EstimatedDays: 4,
		Fee:           0,
	}, nil
}

// FetchStock checks stock availability
func (p *FlipkartProvider) FetchStock(ctx context.Context, url string) (bool, string, error) {
	html, err := p.FetchURL(ctx, url)
	if err != nil {
		return false, "", err
	}

	inStock := !strings.Contains(html, "Currently unavailable") &&
		!strings.Contains(html, "Out of stock") &&
		!strings.Contains(html, "Sold Out")

	status := "in_stock"
	if !inStock {
		status = "out_of_stock"
	}

	return inStock, status, nil
}

// FetchMetadata gets additional product metadata
func (p *FlipkartProvider) FetchMetadata(ctx context.Context, url string) (map[string]interface{}, error) {
	html, err := p.FetchURL(ctx, url)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]interface{})

	// Extract product ID
	if productID, err := p.extractProductID(url); err == nil {
		metadata["product_id"] = productID
	}

	// Extract brand
	brandPattern := regexp.MustCompile(`Brand:\s*([^,<]+)`)
	if matches := brandPattern.FindStringSubmatch(html); len(matches) > 1 {
		metadata["brand"] = strings.TrimSpace(matches[1])
	}

	return metadata, nil
}

// HealthCheck verifies the provider is working
func (p *FlipkartProvider) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := p.FetchURL(ctx, "https://www.flipkart.com")
	return err
}

// Close cleans up resources
func (p *FlipkartProvider) Close() error {
	return nil
}

// extractProductID extracts Flipkart product ID from URL
func (p *FlipkartProvider) extractProductID(url string) (string, error) {
	patterns := []string{
		`/p/([^/]+)`,
		`pid=([A-Z0-9]+)`,
		`([A-Z0-9]{10,})`,
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

// Ensure FlipkartProvider implements base.Provider
var _ base.Provider = (*FlipkartProvider)(nil)
