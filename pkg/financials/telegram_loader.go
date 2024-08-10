package financials

import (
	"fmt"
	"strings"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	databaseQueryErrorFormat string = "Could not query database, error: %s"
	startMonth               int    = 3
	startYear                int    = 2023
)

type TelegramLoader struct {
	fw        framework.FW
	txDB      db.TransactionDB
	expenseDB db.ExpenseDB
}

func NewTelegramLoader(fw framework.FW) Loader {
	return &TelegramLoader{
		fw:        fw,
		txDB:      fw.GetDB(db.TransactionDatabaseName).(db.TransactionDB),
		expenseDB: fw.GetDB(db.ExpenseDatabaseName).(db.ExpenseDB),
	}
}

func (l *TelegramLoader) Load() error {
	if err := l.genExpense(); err != nil {
		return fmt.Errorf("error running genExpense: %s", err)
	}

	return nil
}

func (l *TelegramLoader) genExpense() error {
	loc, _ := time.LoadLocation("Asia/Singapore")
	startMonth, _ := time.ParseInLocation(
		time.DateTime,
		fmt.Sprintf("%d-%02d-01 16:00:00", startYear, startMonth),
		loc,
	)
	today := time.Now()
	for startMonth.Before(today) {
		if err := l.genExpensePerMonth(startMonth); err != nil {
			return err
		}
		if err := l.genSharedExpensePerMonth(startMonth); err != nil {
			return err
		}
		startMonth = startMonth.AddDate(0, 1, 0)
	}

	return nil
}

func (l *TelegramLoader) genExpensePerMonth(startDate time.Time) error {
	endDate := startDate.AddDate(0, 1, 0)

	// Generate TypeOwn (Others field)
	othersResultChannel := make(chan db.AsyncAggregateResult)
	go l.txDB.QueryTypeOwnSum(startDate, endDate, othersResultChannel)

	// Generate Reimbursement (reimbursement field, will be negative)
	reimResultChannel := make(chan db.AsyncAggregateResult)
	go l.txDB.QueryReimburseSum(startDate, endDate, reimResultChannel)

	// Generate shared expenses (list of transactions)
	miscResultChannel := make(chan db.AsyncTransactionResults)
	go l.txDB.QueryMiscTransactions(startDate, endDate, miscResultChannel)

	othersResult := <-othersResultChannel
	if othersResult.Error != nil {
		return fmt.Errorf(databaseQueryErrorFormat, othersResult.Error)
	}

	reimResult := <-reimResultChannel
	if reimResult.Error != nil {
		return fmt.Errorf(databaseQueryErrorFormat, reimResult.Error)
	}

	miscResult := <-miscResultChannel
	if miscResult.Error != nil {
		return fmt.Errorf(databaseQueryErrorFormat, miscResult.Error)
	}

	expenses := []*db.Expense{}
	endOfMonth := utils.DateTime{Time: utils.SetDateToEndOfMonth(startDate)}
	if !utils.AlmostEqual(*othersResult.Result, 0) {
		expenses = append(expenses, &db.Expense{
			TransactionDate: endOfMonth,
			Type:            "Others",
			Amount:          *othersResult.Result,
		})
	}
	if !utils.AlmostEqual(*reimResult.Result, 0) {
		expenses = append(expenses, &db.Expense{
			TransactionDate: endOfMonth,
			Type:            "Reimbursement",
			Amount:          *reimResult.Result * -1,
		})
	}
	for _, result := range miscResult.Result {
		expenses = append(expenses, &db.Expense{
			TransactionDate: endOfMonth,
			Type:            cases.Title(language.English).String(strings.ToLower(string(result.Type))),
			Amount:          result.Amount,
		})
	}

	if len(expenses) > 0 {
		return l.expenseDB.BulkAdd(expenses)
	}
	return nil
}

func (l *TelegramLoader) genSharedExpensePerMonth(startDate time.Time) error {
	endDate := startDate.AddDate(0, 1, 0)

	sharedResultChannel := make(chan db.AsyncTransactionResults)
	go l.txDB.QuerySharedTransactions(startDate, endDate, sharedResultChannel)

	sharedReimCCResultChannel := make(chan db.AsyncTransactionResults)
	go l.txDB.QuerySharedReimCCTransactions(startDate, endDate, sharedReimCCResultChannel)

	sharedResult := <-sharedResultChannel
	if sharedResult.Error != nil {
		return fmt.Errorf(databaseQueryErrorFormat, sharedResult.Error)
	}

	sharedReimCCResult := <-sharedReimCCResultChannel
	if sharedReimCCResult.Error != nil {
		return fmt.Errorf(databaseQueryErrorFormat, sharedReimCCResult.Error)
	}

	nonSpecialSpend := 0.0
	var otherSpendingDate utils.DateTime

	sharedExpenses := []*db.SharedExpense{}
	for _, tx := range sharedResult.Result {
		if strings.Contains(string(tx.Type), "SPECIAL") {
			type_ := fmt.Sprintf("Special:%s", string(tx.Classification))
			endOfMonth := utils.DateTime{Time: utils.SetDateToEndOfMonth(tx.Date)}
			sharedExpenses = append(sharedExpenses, &db.SharedExpense{
				ExpenseDate: endOfMonth,
				Type:        type_,
				Amount:      tx.Amount,
			})
			continue
		}

		// Combine all non special spends
		nonSpecialSpend += tx.Amount
		otherSpendingDate = utils.DateTime{Time: utils.SetDateToEndOfMonth(tx.Date)}
	}

	// subtract shared reim cc result
	sharedCCReimAmount := 0.0
	for _, tx := range sharedReimCCResult.Result {
		sharedCCReimAmount += tx.Amount
	}

	if !utils.AlmostEqual(nonSpecialSpend-sharedCCReimAmount, 0.0) {
		sharedExpenses = append(sharedExpenses, &db.SharedExpense{
			ExpenseDate: otherSpendingDate,
			Type:        "others",
			Amount:      nonSpecialSpend - sharedCCReimAmount,
		})
	}

	if len(sharedExpenses) > 0 {
		return l.expenseDB.BulkAdd(sharedExpenses)
	}

	return nil
}
