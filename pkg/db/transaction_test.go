package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parseDateForced(t *testing.T, dateString string) time.Time {
	loc, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)

	parsed, err := time.ParseInLocation(time.DateOnly, dateString, loc)
	assert.Nil(t, err)

	return parsed
}

func TestInsert(t *testing.T) {
	err := RunTest(func(db *DB) {
		txDB, err := NewTransactionDB(db)
		assert.Nil(t, err)

		tx := &Transaction{
			Date:           time.Now(),
			Type:           TypeReimburse,
			Classification: "negative",
			Amount:         254.23,
		}
		res, err := txDB.InsertTransaction(tx)
		assert.Nil(t, err)
		assert.Equal(t, tx.Type, res.Type)
		assert.Equal(t, tx.Classification, res.Classification)
		assert.Equal(t, tx.Amount, res.Amount)
	})

	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	err := RunTest(func(db *DB) {
		txDB, err := NewTransactionDB(db)
		assert.Nil(t, err)

		tx := &Transaction{
			Date:           time.Now(),
			Type:           TypeReimburse,
			Classification: "negative",
			Amount:         254.23,
		}
		_, err = txDB.InsertTransaction(tx)
		assert.Nil(t, err)

		deletedTx, err := txDB.DeleteTransaction(tx.ID)
		assert.Nil(t, err)

		assert.Equal(t, tx.Type, deletedTx.Type)
		assert.Equal(t, tx.Classification, deletedTx.Classification)
		assert.Equal(t, tx.Amount, deletedTx.Amount)
	})

	assert.Nil(t, err)
}

func TestAggregateTransactions(t *testing.T) {
	var tests = []struct {
		name         string
		inserts      []Transaction
		queryOptions *FindTransactionOptions
		expected     float64
	}{
		{
			name: "get all own",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "meal",
					Amount:         54.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "treat",
					Amount:         46.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeReimburse,
					Classification: "should not be counted",
					Amount:         5.00,
				},
			},
			queryOptions: &FindTransactionOptions{
				StartDate: parseDateForced(t, "2023-04-03"),
				EndDate:   parseDateForced(t, "2023-04-04"),
				Types: []TransactionType{
					TypeOwn,
				},
			},
			expected: 100.46,
		},
		{
			name: "get all reim and shared reim",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeReimburse,
					Classification: "meal",
					Amount:         54.50,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeSharedReimburse,
					Classification: "treat",
					Amount:         46.55,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "should not be counted",
					Amount:         5.00,
				},
			},
			queryOptions: &FindTransactionOptions{
				StartDate: parseDateForced(t, "2023-04-03"),
				EndDate:   parseDateForced(t, "2023-04-04"),
				Types: []TransactionType{
					TypeReimburse,
					TypeSharedReimburse,
				},
			},
			expected: 101.05,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := RunTest(func(db *DB) {
				txDB, err := NewTransactionDB(db)
				assert.Nil(t, err)

				for _, i := range tt.inserts {
					i := i
					_, err := txDB.InsertTransaction(&i)
					assert.Nil(t, err)
				}

				res, err := txDB.AggregateTransactions(tt.queryOptions)
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, *res)

			})
			assert.Nil(t, err)
		})
	}
}

func TestQueryTransactionsByOptions(t *testing.T) {
	var tests = []struct {
		name           string
		inserts        []Transaction
		queryOptions   *FindTransactionOptions
		expectedLength int
	}{
		{
			name: "get all shared reim and shared",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeShared,
					Classification: "meal",
					Amount:         54.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeSharedReimburse,
					Classification: "treat",
					Amount:         46.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "should not be counted",
					Amount:         5.00,
				},
			},
			queryOptions: &FindTransactionOptions{
				StartDate: parseDateForced(t, "2023-04-03"),
				EndDate:   parseDateForced(t, "2023-04-04"),
				Types: []TransactionType{
					TypeSharedReimburse,
					TypeShared,
				},
			},
			expectedLength: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := RunTest(func(db *DB) {
				txDB, err := NewTransactionDB(db)
				assert.Nil(t, err)
				for _, i := range tt.inserts {
					i := i
					_, err := txDB.InsertTransaction(&i)
					assert.Nil(t, err)
				}

				res, err := txDB.QueryTransactionByOptions(tt.queryOptions)
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedLength, len(res))

			})
			assert.Nil(t, err)
		})
	}
}

func TestSearchTransactions(t *testing.T) {
	var tests = []struct {
		name         string
		inserts      []Transaction
		queryOptions *SearchTransactionsOptions
		// expected holds the matching classifications in the order returned.
		expected []string
	}{
		{
			name: "matches the description regardless of case",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "Starbucks Coffee",
					Amount:         5.40,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "COFFEE beans",
					Amount:         12.00,
				},
				{
					Date:           parseDateForced(t, "2023-04-05"),
					Type:           TypeOwn,
					Classification: "lunch",
					Amount:         9.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "coffee",
			},
			// Newest first.
			expected: []string{"COFFEE beans", "Starbucks Coffee"},
		},
		{
			name: "matches a substring of the description",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "lunch at the hawker centre",
					Amount:         4.20,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "dinner",
					Amount:         18.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "HAWKER",
			},
			expected: []string{"lunch at the hawker centre"},
		},
		{
			name: "returns nothing when nothing matches",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "lunch",
					Amount:         4.20,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "nothing like this",
			},
			expected: []string{},
		},
		{
			name: "treats a percent sign literally",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "100% bonus",
					Amount:         100.00,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "100X bonus",
					Amount:         50.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "100%",
			},
			expected: []string{"100% bonus"},
		},
		{
			name: "treats an underscore literally",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "snake_case",
					Amount:         1.00,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "snakecase",
					Amount:         2.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "snake_case",
			},
			expected: []string{"snake_case"},
		},
		{
			name: "filters by type alone",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeSharedCCReimburse,
					Classification: "meal",
					Amount:         54.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "treat",
					Amount:         46.23,
				},
				{
					Date:           parseDateForced(t, "2023-04-05"),
					Type:           TypeTithe,
					Classification: "",
					Amount:         200.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Types: []TransactionType{TypeSharedCCReimburse},
			},
			expected: []string{"meal"},
		},
		{
			name: "filters by description and type together",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "coffee",
					Amount:         5.40,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeShared,
					Classification: "coffee",
					Amount:         6.40,
				},
			},
			queryOptions: &SearchTransactionsOptions{
				Description: "coffee",
				Types:       []TransactionType{TypeOwn},
			},
			expected: []string{"coffee"},
		},
		{
			name: "with no filters returns the most recent transactions",
			inserts: []Transaction{
				{
					Date:           parseDateForced(t, "2023-04-03"),
					Type:           TypeOwn,
					Classification: "oldest",
					Amount:         1.00,
				},
				{
					Date:           parseDateForced(t, "2023-04-05"),
					Type:           TypeOwn,
					Classification: "newest",
					Amount:         2.00,
				},
				{
					Date:           parseDateForced(t, "2023-04-04"),
					Type:           TypeOwn,
					Classification: "middle",
					Amount:         3.00,
				},
			},
			queryOptions: &SearchTransactionsOptions{},
			expected:     []string{"newest", "middle", "oldest"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := RunTest(func(db *DB) {
				txDB, err := NewTransactionDB(db)
				assert.Nil(t, err)
				for _, i := range tt.inserts {
					i := i
					_, err := txDB.InsertTransaction(&i)
					assert.Nil(t, err)
				}

				res, truncated, err := txDB.SearchTransactions(tt.queryOptions)
				assert.Nil(t, err)
				assert.False(t, truncated)

				classifications := make([]string, 0, len(res))
				for _, transaction := range res {
					classifications = append(classifications, transaction.Classification)
				}
				assert.Equal(t, tt.expected, classifications)
			})
			assert.Nil(t, err)
		})
	}
}

func TestSearchTransactionsTruncates(t *testing.T) {
	err := RunTest(func(db *DB) {
		txDB, err := NewTransactionDB(db)
		assert.Nil(t, err)

		// One row more than the cap, so that truncation is detectable.
		inserts := make([]Transaction, 0, MaxTransactionSearchResults+1)
		for i := range MaxTransactionSearchResults + 1 {
			inserts = append(inserts, Transaction{
				Date:           parseDateForced(t, "2023-04-03"),
				Type:           TypeOwn,
				Classification: "match",
				Amount:         float64(i),
			})
		}
		assert.Nil(t, txDB.BulkAdd(inserts))

		res, truncated, err := txDB.SearchTransactions(&SearchTransactionsOptions{
			Description: "match",
		})
		assert.Nil(t, err)
		assert.Equal(t, MaxTransactionSearchResults, len(res))
		assert.True(t, truncated)
	})

	assert.Nil(t, err)
}
