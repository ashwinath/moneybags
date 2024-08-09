package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type Loader struct {
	fw      framework.FW
	assetDB db.AssetDB
}

func NewLoader(fw framework.FW) *Loader {
	assetDB := fw.GetDB(db.AssetDatabaseName).(db.AssetDB)
	return &Loader{
		fw:      fw,
		assetDB: assetDB,
	}
}

func (l *Loader) Start() error {
	dataLoaders := []dataLoader{
		{
			name:     "assets",
			db:       l.assetDB.(db.ClearAndBulkAdder),
			filePath: l.fw.GetConfig().FinancialsData.AssetsCsvFilepath,
			model:    &[]*db.Asset{},
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

func (l *Loader) loadAssets(errChan chan<- error) {
	if err := l.assetDB.Clear(); err != nil {
		errChan <- fmt.Errorf("failed to clear asset db: %s", err)
		return
	}

	assets := []*db.Asset{}
	if err := utils.UnmarshalCSV(l.fw.GetConfig().FinancialsData.AssetsCsvFilepath, &assets); err != nil {
		errChan <- fmt.Errorf("failed to unmarshal csv for assets: %s", err)
		return
	}

	if err := l.assetDB.BulkAdd(assets); err != nil {
		errChan <- fmt.Errorf("failed to add assets: %s", err)
		return
	}

	errChan <- nil
}
