package db

import (
	"fmt"
	"strings"
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
	SearchTransactions(o *SearchTransactionsOptions) ([]Transaction, bool, error)
	QueryTypeOwnSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult)
	QueryReimburseSum(startDate, endDate time.Time, result chan<- AsyncAggregateResult)
	QuerySharedTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	QuerySharedReimCCTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	QueryMiscTransactions(startDate, endDate time.Time, result chan<- AsyncTransactionResults)
	BulkAdd(objs any) error
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

// MaxTransactionSearchResults bounds a single search so that an unbounded
// query, such as one that filters on type alone, cannot return every row.
const MaxTransactionSearchResults = 100

// SearchTransactionsOptions filters a free text search over transaction descriptions.
type SearchTransactionsOptions struct {
	// Description is matched case-insensitively against the transaction's
	// classification. An empty value matches every transaction.
	Description string
	// Types restricts results to the given transaction types. An empty value
	// matches every type.
	Types []TransactionType
}

// SearchTransactions returns the transactions matching o, newest first.
//
// The returned bool reports whether the result cap dropped any further matches,
// so callers can tell a short result set apart from a truncated one. At most
// MaxTransactionSearchResults transactions are returned.
func (d *transactionDB) SearchTransactions(o *SearchTransactionsOptions) ([]Transaction, bool, error) {
	query := d.db.Model(&Transaction{})

	if o.Description != "" {
		// ILIKE performs the comparison case-insensitively in the database. The
		// user's text is escaped so that a literal % or _ in a description is
		// matched instead of being read as a wildcard.
		query = query.Where(
			`classification ILIKE ? ESCAPE '\'`,
			"%"+escapeLikePattern(o.Description)+"%",
		)
	}

	if len(o.Types) > 0 {
		query = query.Where("type in ?", o.Types)
	}

	// One row beyond the cap is read so that truncation can be reported
	// accurately rather than guessed at from a full page of results.
	var transactions []Transaction
	result := query.
		Order("date desc").
		Limit(MaxTransactionSearchResults + 1).
		Find(&transactions)
	if result.Error != nil {
		return nil, false, result.Error
	}

	truncated := len(transactions) > MaxTransactionSearchResults
	if truncated {
		transactions = transactions[:MaxTransactionSearchResults]
	}

	return transactions, truncated, nil
}

// escapeLikePattern escapes the wildcards that LIKE and ILIKE treat specially,
// so that they are matched literally.
func escapeLikePattern(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
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
func (db *transactionDB) BulkAdd(objs any) error {
	return db.db.Create(objs).Error
}
