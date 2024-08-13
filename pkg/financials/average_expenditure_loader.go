package financials

import (
	"fmt"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

const (
	windowPeriod = 6
)

type averageExpenditureLoader struct {
	fw                   framework.FW
	averageExpenditureDB db.AverageExpenditureDB
	expenseDB            db.ExpenseDB
}

func NewAverageExpenditureLoader(fw framework.FW) Loader {
	return &averageExpenditureLoader{
		fw:                   fw,
		averageExpenditureDB: fw.GetDB(db.AverageExpenditureDatabaseName).(db.AverageExpenditureDB),
		expenseDB:            fw.GetDB(db.ExpenseDatabaseName).(db.ExpenseDB),
	}
}

func (averageExpenditureLoader) Name() string {
	return "average expenditure"
}

func (l *averageExpenditureLoader) Load() error {
	firstDate, err := l.expenseDB.GetFirstDate()
	if err != nil {
		return fmt.Errorf("failed to get first date of portfolio")
	}

	currentDate := firstDate
	if currentDate.Day() != 1 {
		currentDate = time.Date(currentDate.Year(), currentDate.Month(), 1, currentDate.Hour(), 0, 0, 0, currentDate.Location())
		currentDate = currentDate.AddDate(0, 1, 0)
	}

	allAverageExpenditures := []db.AverageExpenditure{}
	tomorrow := time.Now().AddDate(0, 0, 1)

	for currentDate.Before(tomorrow) {
		yearlyExpenditure, err := l.expenseDB.GetYearlyExpense([]string{"Tax", "Special:%"}, currentDate, windowPeriod)
		if err != nil {
			// That period could be empty.
			currentDate = utils.GetLastDateOfMonth(currentDate.AddDate(0, 1, 0))
			continue
		}

		ae := db.AverageExpenditure{
			ExpenseDate: currentDate,
			Amount:      yearlyExpenditure / windowPeriod,
		}
		allAverageExpenditures = append(allAverageExpenditures, ae)
		currentDate = utils.GetLastDateOfMonth(currentDate.AddDate(0, 1, 0))
	}

	if err := l.averageExpenditureDB.BulkInsertOnConflictOverride(allAverageExpenditures); err != nil {
		return fmt.Errorf("failed to bulk insert average expenditures: %s", err)
	}
	return nil
}
