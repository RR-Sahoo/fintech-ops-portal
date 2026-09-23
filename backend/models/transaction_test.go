package models_test

import (
	"testing"

	"github.com/RR-Sahoo/fintech-starter-kit/backend/models"
)

func TestFormatUSD(t *testing.T) {
	tests := []struct {
		name     string
		cents    int64
		expected string
	}{
		{"Zero cents", 0, "$0.00"},
		{"50 cents", 50, "$0.50"},
		{"One dollar", 100, "$1.00"},
		{"One hundred twenty five dollars", 12500, "$125.00"},
		{"Five hundred dollars", 50000, "$500.00"},
		{"One thousand two hundred fifty dollars", 125000, "$1,250.00"},
		{"Four thousand five hundred dollars", 450000, "$4,500.00"},
		{"Negative amount", -125000, "-$1,250.00"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := models.FormatUSD(tc.cents)
			if got != tc.expected {
				t.Errorf("FormatUSD(%d) = %q; want %q", tc.cents, got, tc.expected)
			}
		})
	}
}
