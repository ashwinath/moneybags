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
