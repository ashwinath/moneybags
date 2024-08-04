package db

import (
	"fmt"
	"time"
)

type TransactionType string

const TypeReimburse TransactionType = "REIM"
const TypeSharedReimburse TransactionType = "SHARED_REIM"
const TypeShared TransactionType = "SHARED"
const TypeSpecialShared TransactionType = "SPECIAL_SHARED"
const TypeSpecialSharedReimburse TransactionType = "SPECIAL_SHARED_REIM"
const TypeOwn TransactionType = "OWN"
const TypeSpecialOwn TransactionType = "SPECIAL_OWN"
const TypeCreditCard TransactionType = "CREDIT CARD"
const TypeInsurance TransactionType = "INSURANCE"
const TypeTithe TransactionType = "TITHE"
const TypeTax TransactionType = "TAX"
const TypeSharedCCReimburse TransactionType = "SHARED_CC_REIMBURSE"

type Transaction struct {
	ID             uint            `gorm:"primaryKey"`
	Date           time.Time       `gorm:"type:timestamp;index"`
	Type           TransactionType `gorm:"column:type"`
	Classification string
	Amount         float64
	CreatedAt      time.Time `gorm:"type:timestamp"`
	UpdatedAt      time.Time `gorm:"type:timestamp"`
}

func (t *Transaction) String() string {
	if len(string(t.Classification)) == 0 {
		return fmt.Sprintf(
			"Date: %s\nType: %s\nAmount: %.2f",
			t.Date,
			t.Type,
			t.Amount,
		)
	}

	return fmt.Sprintf(
		"Date: %s\nType: %s\nClassification: %s\nAmount:%.2f",
		t.Date,
		t.Type,
		t.Classification,
		t.Amount,
	)
}

// TODO: define functions
type TransactionDB interface{}

type transactionDB struct {
	*DB
}

func NewTransactionDB(db *DB) (TransactionDB, error) {
	if err := db.DB.AutoMigrate(&Transaction{}); err != nil {
		return nil, err
	}

	return &transactionDB{
		DB: db,
	}, nil
}
