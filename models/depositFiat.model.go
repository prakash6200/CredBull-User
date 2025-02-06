package models

import (
	"time"

	"gorm.io/gorm"
)

// Status Enum
type Status string

const (
	Pending    Status = "pending"
	Approved   Status = "approved"
	Rejected   Status = "rejected"
	Processing Status = "processing"
)

// BankPaymentType Enum
type BankPaymentType string

const (
	NEFT BankPaymentType = "NEFT"
	IMPS BankPaymentType = "IMPS"
	UPI  BankPaymentType = "UPI"
)

type DepositFiat struct {
	gorm.Model
	UserID          string          `gorm:"not null" json:"user_id"`
	Amount          float64         `gorm:"default:0" json:"fiat_amount"`
	Currency        string          `json:"currency"`
	TransactionID   string          `json:"transaction_id"`
	Status          Status          `gorm:"type:varchar(20);default:'pending';check:status IN ('pending', 'approved', 'rejected', 'processing')" json:"status"`
	BankPaymentType BankPaymentType `gorm:"type:varchar(10);check:bank_payment_type IN ('NEFT', 'IMPS', 'UPI')" json:"bank_payment_type"`
	Image           string          `gorm:"default:''" json:"image"`
	ApprovedBy      string          `json:"approved_by"`
	Description     string          `json:"description"`
	AdminActionDate *time.Time      `json:"admin_action_date"`
}
