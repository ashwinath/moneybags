package financials

import (
	"fmt"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type houseAssetLoader struct {
	fw         framework.FW
	assetsDB   db.AssetDB
	mortgageDB db.MortgageDB
}

func NewHouseAssetLoader(fw framework.FW) Loader {
	return &houseAssetLoader{
		fw:         fw,
		assetsDB:   fw.GetDB(db.AssetDatabaseName).(db.AssetDB),
		mortgageDB: fw.GetDB(db.MortgageDatabaseName).(db.MortgageDB),
	}
}

func (houseAssetLoader) Name() string {
	return "house asset"
}

func (l *houseAssetLoader) Load() error {
	mortgages, err := l.mortgageDB.GetMortgage()
	if err != nil {
		return fmt.Errorf("unable to query mortgage: %s", err)
	}

	assets := []db.Asset{}
	for _, m := range mortgages {
		asset := db.Asset{
			TransactionDate: utils.DateTime{Time: m.Date},
			Type:            "House",
			Amount:          m.TotalInterestPaid / numberOfPeopleSharing,
		}
		assets = append(assets, asset)
	}

	if len(assets) == 0 {
		return nil
	}

	allHouseAssets := []db.Asset{}
	// Find gaps in between dates that have no mortgage schedule
	for i := 0; i < len(assets)-1; i++ {
		asset := assets[i]
		currentDate := asset.TransactionDate.Time
		nextAsset := assets[i+1]

		for currentDate.Before(nextAsset.TransactionDate.Time) {
			injectedAsset := db.Asset{
				TransactionDate: utils.DateTime{Time: currentDate},
				Type:            "House",
				Amount:          asset.Amount,
			}
			currentDate = currentDate.AddDate(0, 1, 0)
			allHouseAssets = append(allHouseAssets, injectedAsset)
		}
	}

	if err := l.assetsDB.BulkAdd(allHouseAssets); err != nil {
		return fmt.Errorf("failed to bulk insert house data into asset db: %s", err)
	}

	return nil
}
