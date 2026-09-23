package models_test

import (
	"testing"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
)

func TestFormatRupees(t *testing.T) {
	tests := []struct {
		name     string
		paise    int64
		expected string
	}{
		{"Zero paise", 0, "₹0.00"},
		{"50 paise", 50, "₹0.50"},
		{"One rupee", 100, "₹1.00"},
		{"One thousand two hundred fifty rupees and 50 paise", 125050, "₹1,250.50"},
		{"Fifty thousand rupees", 5000000, "₹50,000.00"},
		{"One lakh twenty five thousand rupees", 12500000, "₹1,25,000.00"},
		{"Negative amount", -125050, "-₹1,250.50"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := models.FormatRupees(tc.paise)
			if got != tc.expected {
				t.Errorf("FormatRupees(%d) = %q; want %q", tc.paise, got, tc.expected)
			}
		})
	}
}
