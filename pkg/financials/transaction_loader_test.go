package financials

import (
	"fmt"
	"testing"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestTransactionLoaderEmpty(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)

		// Test
		loader := NewTransactionLoader(fw)
		err := loader.Load()
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

func TestTransactionLoader(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)

		// Add some dummy data for telegram
		transactions := []*database.Transaction{
			{
				Date:           parseDateForced(t, "2023-04-03"),
				Type:           database.TypeOwn,
				Classification: "meal",
				Amount:         54.23,
			},
			{
				Date:           parseDateForced(t, "2023-04-03"),
				Type:           database.TypeOwn,
				Classification: "treat",
				Amount:         46.23,
			},
			{
				Date:           parseDateForced(t, "2023-04-03"),
				Type:           database.TypeReimburse,
				Classification: "bought ckt for others",
				Amount:         5.00,
			},
			{
				Date:           parseDateForced(t, "2023-04-03"),
				Type:           database.TypeSharedReimburse,
				Classification: "bought shared ckt with credit card",
				Amount:         5.00,
			},
		}
		txDB := fw.GetDB(database.TransactionDatabaseName).(database.TransactionDB)
		err := txDB.BulkAdd(&transactions)
		assert.Nil(t, err)

		// Test
		loader := NewTransactionLoader(fw)
		err = loader.Load()
		assert.Nil(t, err)

		// Assert shared
		sharedExpensesResult := []database.SharedExpense{}
		res := db.DB.Find(&sharedExpensesResult)
		assert.Nil(t, res.Error)
		assert.Len(t, sharedExpensesResult, 1)
		assert.Equal(t, sharedExpensesResult[0].Amount, float64(5))

		// assert expenses
		expensesResult := []database.Expense{}
		res = db.DB.Find(&expensesResult)
		assert.Nil(t, res.Error)
		assert.Len(t, expensesResult, 2)
		for _, r := range expensesResult {
			if r.Type == "Others" {
				assert.Equal(t, r.Amount, float64(100.46))
			} else if r.Type == "Reimbursement" {
				assert.Equal(t, r.Amount, float64(-10))
			} else {
				assert.Error(t, fmt.Errorf("there should be no other case"))
			}
		}
	})
	assert.Nil(t, err)
}
