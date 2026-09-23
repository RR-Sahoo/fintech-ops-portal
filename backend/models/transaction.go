package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Transaction status constants
const (
	StatusSuccess = "SUCCESS"
	StatusPending = "PENDING"
	StatusFailed  = "FAILED"
)

// Transaction type constants
const (
	TypeDebit  = "DEBIT"
	TypeCredit = "CREDIT"
)

// Transaction represents a financial transaction model.
// AmountCents strictly uses int64 to prevent floating-point precision loss.
type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	AmountCents     int64     `json:"amount_cents"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	ReferenceID     string    `json:"reference_id"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	FormattedAmount string    `json:"formatted_amount"`
}

// FormatUSD converts an amount in cents (int64) into a USD formatted string (e.g. $1,250.50).
// Pure integer arithmetic without floating-point precision loss.
func FormatUSD(cents int64) string {
	isNegative := cents < 0
	if isNegative {
		cents = -cents
	}
	dollars := cents / 100
	remainderCents := cents % 100

	dollarsStr := strconv.FormatInt(dollars, 10)
	var formattedDollars strings.Builder
	n := len(dollarsStr)

	// Standard thousands grouping for USD: 1,234,567
	if n <= 3 {
		formattedDollars.WriteString(dollarsStr)
	} else {
		pre := n % 3
		if pre == 0 {
			pre = 3
		}
		formattedDollars.WriteString(dollarsStr[:pre])
		for i := pre; i < n; i += 3 {
			formattedDollars.WriteString(",")
			formattedDollars.WriteString(dollarsStr[i : i+3])
		}
	}

	prefix := "$"
	if isNegative {
		prefix = "-$"
	}
	return fmt.Sprintf("%s%s.%02d", prefix, formattedDollars.String(), remainderCents)
}

// GetFormattedAmount returns the transaction's formatted amount in USD.
func (t *Transaction) GetFormattedAmount() string {
	return FormatUSD(t.AmountCents)
}
