package financials

import (
	"fmt"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

type investmentsLoader struct {
	fw          framework.FW
	assetsDB    db.AssetDB
	portfolioDB db.PortfolioDB
}

func NewInvestmentsLoader(fw framework.FW) Loader {
	return &investmentsLoader{
		fw:          fw,
		assetsDB:    fw.GetDB(db.AssetDatabaseName).(db.AssetDB),
		portfolioDB: fw.GetDB(db.PortfolioDatabaseName).(db.PortfolioDB),
	}
}

func (l *investmentsLoader) Load() error {
	firstDate, err := l.portfolioDB.GetFirstTradeDate()
	if err != nil {
		return fmt.Errorf("failed to get first date of portfolio")
	}

	currentDate := firstDate
	if currentDate.Day() != 1 {
		currentDate = time.Date(currentDate.Year(), currentDate.Month(), 1, currentDate.Hour(), 0, 0, 0, currentDate.Location())
		currentDate = currentDate.AddDate(0, 1, 0)
	}

	allInvestments := []db.Asset{}
	tomorrow := time.Now().AddDate(0, 0, 1)

	for currentDate.Before(tomorrow) {
		amount, err := l.portfolioDB.GetPortfolioAmountByDate(currentDate)
		if err != nil {
			return fmt.Errorf("failed to get amount for portfolio by date (%s): %s", currentDate, err)
		}

		asset := db.Asset{
			TransactionDate: utils.DateTime{Time: currentDate},
			Type:            "Investments",
			Amount:          amount,
		}
		allInvestments = append(allInvestments, asset)

		currentDate = currentDate.AddDate(0, 1, 0)
	}

	if err := l.assetsDB.BulkAdd(allInvestments); err != nil {
		return fmt.Errorf("failed to bulk insert investments into assets db: %s", err)
	}

	return nil
}
