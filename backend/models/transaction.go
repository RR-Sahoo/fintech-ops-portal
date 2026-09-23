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
// AmountPaise strictly uses int64 to prevent floating-point precision loss.
type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	AmountPaise     int64     `json:"amount_paise"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	ReferenceID     string    `json:"reference_id"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	FormattedAmount string    `json:"formatted_amount"`
}

// FormatRupees converts an amount in paise (int64) into an Indian Rupee formatted string (e.g. ₹1,250.50).
// Does not use any floating point operations.
func FormatRupees(paise int64) string {
	isNegative := paise < 0
	if isNegative {
		paise = -paise
	}
	rupees := paise / 100
	remainderPaise := paise % 100

	rupeesStr := strconv.FormatInt(rupees, 10)
	var formattedRupees strings.Builder
	n := len(rupeesStr)

	if n <= 3 {
		formattedRupees.WriteString(rupeesStr)
	} else {
		last3 := rupeesStr[n-3:]
		remaining := rupeesStr[:n-3]

		var parts []string
		for len(remaining) > 2 {
			parts = append([]string{remaining[len(remaining)-2:]}, parts...)
			remaining = remaining[:len(remaining)-2]
		}
		if len(remaining) > 0 {
			parts = append([]string{remaining}, parts...)
		}
		formattedRupees.WriteString(strings.Join(parts, ","))
		formattedRupees.WriteString(",")
		formattedRupees.WriteString(last3)
	}

	prefix := "₹"
	if isNegative {
		prefix = "-₹"
	}
	return fmt.Sprintf("%s%s.%02d", prefix, formattedRupees.String(), remainderPaise)
}

// FormattedAmount returns the transaction's formatted amount in rupees.
func (t *Transaction) GetFormattedAmount() string {
	return FormatRupees(t.AmountPaise)
}
