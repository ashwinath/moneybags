package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type Loader struct {
	fw              framework.FW
	assetDB         db.ClearAndBulkAdder
	expenseDB       db.ClearAndBulkAdder
	incomeDB        db.ClearAndBulkAdder
	sharedExpenseDB db.ClearAndBulkAdder
	tradesDB        db.ClearAndBulkAdder
}

func NewLoader(fw framework.FW) *Loader {
	return &Loader{
		fw:              fw,
		assetDB:         fw.GetDB(db.AssetDatabaseName).(db.ClearAndBulkAdder),
		expenseDB:       fw.GetDB(db.ExpenseDatabaseName).(db.ClearAndBulkAdder),
		incomeDB:        fw.GetDB(db.IncomeDatabaseName).(db.ClearAndBulkAdder),
		sharedExpenseDB: fw.GetDB(db.IncomeDatabaseName).(db.ClearAndBulkAdder),
		tradesDB:        fw.GetDB(db.TradeDatabaseName).(db.ClearAndBulkAdder),
	}
}

func (l *Loader) Start() error {
	dataLoaders := []dataLoader{
		{
			name:     "assets",
			db:       l.assetDB,
			filePath: l.fw.GetConfig().FinancialsData.AssetsCsvFilepath,
			model:    &[]*db.Asset{},
			errChan:  make(chan error, 1),
		},
		{
			name:     "expenses",
			db:       l.expenseDB,
			filePath: l.fw.GetConfig().FinancialsData.ExpensesCsvFilepath,
			model:    &[]*db.Expense{},
			errChan:  make(chan error, 1),
		},
		{
			name:     "incomeDB",
			db:       l.incomeDB,
			filePath: l.fw.GetConfig().FinancialsData.IncomeCsvFilepath,
			model:    &[]*db.Income{},
			errChan:  make(chan error, 1),
		},
		{
			name:     "sharedExpenseDB",
			db:       l.sharedExpenseDB,
			filePath: l.fw.GetConfig().FinancialsData.SharedExpensesCsvFilepath,
			model:    &[]*db.SharedExpense{},
			errChan:  make(chan error, 1),
		},
		{
			name:     "trades",
			db:       l.tradesDB,
			filePath: l.fw.GetConfig().FinancialsData.TradesCsvFilepath,
			model:    &[]*db.Trade{},
			errChan:  make(chan error, 1),
		},
	}

	for _, d := range dataLoaders {
		go d.load()
	}

	for _, d := range dataLoaders {
		if err := <-d.errChan; err != nil {
			return err
		}
	}

	return nil
}

type dataLoader struct {
	name     string
	db       db.ClearAndBulkAdder
	filePath string
	model    interface{}
	errChan  chan error
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
