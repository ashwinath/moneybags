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
	if err := l.loadAssets(); err != nil {
		return fmt.Errorf("failed to load assets: %s", err)
	}

	return nil
}

func (l *Loader) loadAssets() error {
	if err := l.assetDB.Clear(); err != nil {
		return fmt.Errorf("failed to clear asset db: %s", err)
	}

	assets := []*db.Asset{}
	if err := utils.UnmarshalCSV(l.fw.GetConfig().FinancialsData.AssetsCsvFilepath, &assets); err != nil {
		return fmt.Errorf("failed to unmarshal csv for assets: %s", err)
	}

	if err := l.assetDB.BulkAdd(assets); err != nil {
		return fmt.Errorf("failed to add assets: %s", err)
	}

	return nil
}
