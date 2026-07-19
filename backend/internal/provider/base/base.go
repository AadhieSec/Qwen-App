package base

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// BaseProvider provides common functionality for all providers
type BaseProvider struct {
	HTTPClient *http.Client
	UserAgent  string
	Timeout    time.Duration
}

// NewBaseProvider creates a new base provider with default settings
func NewBaseProvider(userAgent string, timeout time.Duration) *BaseProvider {
	return &BaseProvider{
		HTTPClient: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				return nil
			},
		},
		UserAgent: userAgent,
		Timeout:   timeout,
	}
}

// FetchURL fetches a URL and returns the response body
func (b *BaseProvider) FetchURL(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", b.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	resp, err := b.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	buf := make([]byte, 1024*1024) // 1MB buffer
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	return string(buf[:n]), nil
}

// ParseHTML parses HTML content into a goquery document
func (b *BaseProvider) ParseHTML(html string) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(nil)
}

// ExtractPrice extracts price from a string
func (b *BaseProvider) ExtractPrice(priceStr string) float64 {
	// Remove currency symbols and commas
	cleanStr := priceStr
	replacements := []struct {
		old, new string
	}{
		{"₹", ""},
		{"$", ""},
		{"€", ""},
		{"£", ""},
		{",", ""},
		{" ", ""},
	}

	for _, r := range replacements {
		cleanStr = replaceAll(cleanStr, r.old, r.new)
	}

	var price float64
	fmt.Sscanf(cleanStr, "%f", &price)
	return price
}

// CalculateDiscount calculates discount percentage
func (b *BaseProvider) CalculateDiscount(price, mrp float64) float64 {
	if mrp <= 0 {
		return 0
	}
	return ((mrp - price) / mrp) * 100
}

// DetermineStockStatus determines stock status from various indicators
func (b *BaseProvider) DetermineStockStatus(inStock bool, statusText string) string {
	if inStock {
		return "in_stock"
	}
	if statusText != "" {
		return statusText
	}
	return "out_of_stock"
}

// Helper function to replace all occurrences
func replaceAll(s, old, new string) string {
	result := s
	for i := 0; i < len(result); i++ {
		if i+len(old) <= len(result) && result[i:i+len(old)] == old {
			result = result[:i] + new + result[i+len(old):]
			i += len(new) - 1
		}
	}
	return result
}
