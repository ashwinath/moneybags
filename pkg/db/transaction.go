package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const TransactionDatabaseName string = "transaction"

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
	Date           time.Time       `gorm:"type:timestamptz;index"`
	Type           TransactionType `gorm:"column:type"`
	Classification string
	Amount         float64
	CreatedAt      time.Time `gorm:"type:timestamptz"`
	UpdatedAt      time.Time `gorm:"type:timestamptz"`
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
		"ID: %d\nDate: %s\nType: %s\nClassification: %s\nAmount:%.2f",
		t.ID,
		t.Date,
		t.Type,
		t.Classification,
		t.Amount,
	)
}

type TransactionDB interface {
	InsertTransaction(tx *Transaction) (*Transaction, error)
	DeleteTransaction(id uint) (*Transaction, error)
	AggregateTransactions(o *FindTransactionOptions) (*float64, error)
	QueryTransactionByOptions(o *FindTransactionOptions) ([]Transaction, error)
	QueryTypeOwnSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult)
	QueryReimburseSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult)
	QuerySharedTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	QuerySharedReimCCTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	QueryMiscTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	BulkAdd(objs interface{}) error
}

type transactionDB struct {
	db *gorm.DB
}

func NewTransactionDB(db *DB) (TransactionDB, error) {
	if err := db.DB.AutoMigrate(&Transaction{}); err != nil {
		return nil, err
	}

	return &transactionDB{
		db: db.DB,
	}, nil
}

func (d *transactionDB) InsertTransaction(tx *Transaction) (*Transaction, error) {
	result := d.db.Create(tx)
	if result.Error != nil {
		return nil, result.Error
	}

	return d.queryTransaction(tx.ID)
}

func (d *transactionDB) queryTransaction(id uint) (*Transaction, error) {
	tx := &Transaction{}
	result := d.db.First(tx, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return tx, nil
}

// returns the old copy of the deleted transaction
func (d *transactionDB) DeleteTransaction(id uint) (*Transaction, error) {
	deletedTx, err := d.queryTransaction(id)
	if err != nil {
		return nil, err
	}

	result := d.db.Delete(&Transaction{ID: id})
	if result.Error != nil {
		return nil, result.Error
	}

	return deletedTx, nil
}

type FindTransactionOptions struct {
	StartDate time.Time
	EndDate   time.Time
	Types     []TransactionType
}

type findTransactionResult struct {
	Total float64
}

func (d *transactionDB) AggregateTransactions(o *FindTransactionOptions) (*float64, error) {
	result := findTransactionResult{}
	res := d.db.Model(&Transaction{}).
		Select("sum(amount) as total").
		Where("date >= ? and date < ? and type in ?", o.StartDate, o.EndDate, o.Types).
		Scan(&result)

	if res.Error != nil {
		return nil, res.Error
	}

	return &result.Total, nil
}

func (d *transactionDB) QueryTransactionByOptions(o *FindTransactionOptions) ([]Transaction, error) {
	var transactions []Transaction
	result := d.db.
		Where("date >= ? and date < ? and type in ?", o.StartDate, o.EndDate, o.Types).
		Order("date asc").
		Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

type AsyncAggregateResult struct {
	Result *float64
	Error  error
}

type AsyncTransactionResults struct {
	Result []Transaction
	Error  error
}

func (d *transactionDB) QueryTypeOwnSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult) {
	defer close(result)
	othersOption := &FindTransactionOptions{
		StartDate: startDate,
		EndDate:   endDate,
		Types:     []TransactionType{TypeOwn},
	}

	othersTotal, err := d.AggregateTransactions(othersOption)
	result <- AsyncAggregateResult{
		Result: othersTotal,
		Error:  err,
	}
}

func (d *transactionDB) QueryReimburseSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult) {
	defer close(result)
	reimOption := &FindTransactionOptions{
		StartDate: startDate,
		EndDate:   endDate,
		Types: []TransactionType{
			TypeReimburse,
			TypeSharedReimburse,
			TypeSpecialSharedReimburse,
		},
	}

	reimTotal, err := d.AggregateTransactions(reimOption)
	result <- AsyncAggregateResult{
		Result: reimTotal,
		Error:  err,
	}
}

func (d *transactionDB) QuerySharedTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults) {
	defer close(result)
	sharedOption := &FindTransactionOptions{
		StartDate: startDate,
		EndDate:   endDate,
		Types: []TransactionType{
			TypeSharedReimburse,
			TypeSpecialSharedReimburse,
			TypeSpecialShared,
			TypeShared,
		},
	}
	sharedTransactions, err := d.QueryTransactionByOptions(sharedOption)
	result <- AsyncTransactionResults{
		Result: sharedTransactions,
		Error:  err,
	}
}

func (d *transactionDB) QuerySharedReimCCTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults) {
	defer close(result)
	sharedOption := &FindTransactionOptions{
		StartDate: startDate,
		EndDate:   endDate,
		Types: []TransactionType{
			TypeSharedCCReimburse,
		},
	}
	sharedTransactions, err := d.QueryTransactionByOptions(sharedOption)
	result <- AsyncTransactionResults{
		Result: sharedTransactions,
		Error:  err,
	}
}

func (d *transactionDB) QueryMiscTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults) {
	defer close(result)
	sharedOption := &FindTransactionOptions{
		StartDate: startDate,
		EndDate:   endDate,
		Types: []TransactionType{
			TypeCreditCard,
			TypeInsurance,
			TypeTax,
			TypeTithe,
		},
	}
	sharedTransactions, err := d.QueryTransactionByOptions(sharedOption)
	result <- AsyncTransactionResults{
		Result: sharedTransactions,
		Error:  err,
	}
}

// Bulk add data
func (db *transactionDB) BulkAdd(objs interface{}) error {
	return db.db.Create(objs).Error
}
