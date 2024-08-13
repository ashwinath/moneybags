package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type csvLoader struct {
	fw      framework.FW
	loaders []dataLoader
}

type dataLoader struct {
	name     string
	db       db.ClearAndBulkAdder
	filePath string
	model    interface{}
	errChan  chan error
}

func NewCSVLoader(fw framework.FW) Loader {
	return &csvLoader{
		fw: fw,
		loaders: []dataLoader{
			{
				name:     "assets",
				db:       fw.GetDB(db.AssetDatabaseName).(db.ClearAndBulkAdder),
				filePath: fw.GetConfig().FinancialsData.AssetsCsvFilepath,
				model:    &[]*db.Asset{},
				errChan:  make(chan error, 1),
			},
			{
				name:     "expenses",
				db:       fw.GetDB(db.ExpenseDatabaseName).(db.ClearAndBulkAdder),
				filePath: fw.GetConfig().FinancialsData.ExpensesCsvFilepath,
				model:    &[]*db.Expense{},
				errChan:  make(chan error, 1),
			},
			{
				name:     "incomeDB",
				db:       fw.GetDB(db.IncomeDatabaseName).(db.ClearAndBulkAdder),
				filePath: fw.GetConfig().FinancialsData.IncomeCsvFilepath,
				model:    &[]*db.Income{},
				errChan:  make(chan error, 1),
			},
			{
				name:     "sharedExpenseDB",
				db:       fw.GetDB(db.SharedExpenseDatabaseName).(db.ClearAndBulkAdder),
				filePath: fw.GetConfig().FinancialsData.SharedExpensesCsvFilepath,
				model:    &[]*db.SharedExpense{},
				errChan:  make(chan error, 1),
			},
			{
				name:     "trades",
				db:       fw.GetDB(db.TradeDatabaseName).(db.ClearAndBulkAdder),
				filePath: fw.GetConfig().FinancialsData.TradesCsvFilepath,
				model:    &[]*db.Trade{},
				errChan:  make(chan error, 1),
			},
		},
	}
}

func (csvLoader) Name() string {
	return "csv"
}

func (l *csvLoader) Load() error {
	for _, d := range l.loaders {
		go d.load()
	}

	for _, d := range l.loaders {
		if err := <-d.errChan; err != nil {
			return err
		}
	}

	return nil
}

func (a *dataLoader) load() {
	if err := a.db.Clear(); err != nil {
		a.errChan <- fmt.Errorf("failed to clear %s db: %s", a.name, err)
		return
	}

	if err := utils.UnmarshalCSV(a.filePath, a.model); err != nil {
		a.errChan <- fmt.Errorf("failed to unmarshal csv for %s: %s", a.name, err)
		return
	}

	if err := a.db.BulkAdd(a.model); err != nil {
		a.errChan <- fmt.Errorf("failed to add %s: %s", a.name, err)
		return
	}

	a.errChan <- nil
}
