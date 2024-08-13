package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

const (
	nonSpecialSharedExpensesType = "Shared Expense"
	specialSharedExpensesType    = "Special:Shared Expense"
	numberOfPeopleSharing        = 2.0
)

type sharedExpenseLoader struct {
	fw              framework.FW
	sharedExpenseDB db.SharedExpenseDB
	expenseDB       db.ExpenseDB
}

func NewSharedExpenseLoader(fw framework.FW) Loader {
	return &sharedExpenseLoader{
		fw:              fw,
		expenseDB:       fw.GetDB(db.ExpenseDatabaseName).(db.ExpenseDB),
		sharedExpenseDB: fw.GetDB(db.SharedExpenseDatabaseName).(db.SharedExpenseDB),
	}
}

func (sharedExpenseLoader) Name() string {
	return "shared expense"
}

func (l *sharedExpenseLoader) Load() error {
	nonSpecialSharedExpenses, err := l.sharedExpenseDB.GetSharedExpensesGroupByExpenseDate(false)
	if err != nil {
		return fmt.Errorf("failed to query non special shared expenses: %s", err)
	}

	allExpenses := []db.Expense{}
	for _, e := range nonSpecialSharedExpenses {
		expense := db.Expense{
			TransactionDate: utils.DateTime{Time: e.ExpenseDate},
			Type:            nonSpecialSharedExpensesType,
			Amount:          e.Total / numberOfPeopleSharing,
		}
		allExpenses = append(allExpenses, expense)
	}

	specialSharedExpenses, err := l.sharedExpenseDB.GetSharedExpensesGroupByExpenseDate(true)
	if err != nil {
		return fmt.Errorf("failed to query non special shared expenses: %s", err)
	}

	for _, e := range specialSharedExpenses {
		expense := db.Expense{
			TransactionDate: utils.DateTime{Time: e.ExpenseDate},
			Type:            specialSharedExpensesType,
			Amount:          e.Total / numberOfPeopleSharing,
		}
		allExpenses = append(allExpenses, expense)
	}

	if err := l.expenseDB.BulkAdd(allExpenses); err != nil {
		return fmt.Errorf("failed to insert special expenses into expense db: %s", err)
	}

	return nil
}
