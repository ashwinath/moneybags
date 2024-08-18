package financials

import (
	"fmt"
	"math"
	"time"

	"github.com/ashwinath/moneybags/pbgo/carpb"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

const (
	assetTypeCarFormat string = "Car:%s"
)

type carLoader struct {
	fw        framework.FW
	assetsDB  db.AssetDB
	carConfig *carpb.CarConfig
}

func NewCarLoader(fw framework.FW) Loader {
	return &carLoader{
		fw:       fw,
		assetsDB: fw.GetDB(db.AssetDatabaseName).(db.AssetDB),
	}
}

func (carLoader) Name() string {
	return "car"
}

func (l *carLoader) Load() error {
	if err := l.loadCarConfig(); err != nil {
		return err
	}

	carAssets := []db.Asset{}
	for _, car := range l.carConfig.Cars {
		a, err := l.GetAssetsPerCar(car)
		if err != nil {
			return fmt.Errorf("could not process car (%s): %s", car.Name, err)
		}
		carAssets = append(carAssets, a...)
	}

	if err := l.assetsDB.BulkAdd(carAssets); err != nil {
		return fmt.Errorf("failed to bulk insert car assets into assets db: %s", err)
	}

	return nil
}

func (l *carLoader) loadCarConfig() error {
	carConfig := carpb.CarConfig{}
	if err := utils.UnmarshalYAML(l.fw.GetConfig().FinancialsData.CarYamlFilepath, &carConfig); err != nil {
		return fmt.Errorf("failed to unmarshal car config: %s", err)
	}

	l.carConfig = &carConfig
	return nil
}

// public for unit test
func (l *carLoader) GetAssetsPerCar(car *carpb.Car) ([]db.Asset, error) {
	carStartDate, err := utils.SetDateFromString(car.CarStartDate)
	if err != nil {
		return nil, fmt.Errorf("could not parse start date (%s): %s", car.CarStartDate, err)
	}
	var carSoldDate time.Time
	if car.CarSoldDate != "" {
		carSoldDate, err = utils.SetDateFromString(car.CarSoldDate)
		if err != nil {
			return nil, fmt.Errorf("could not parse end date (%s): %s", car.CarSoldDate, err)
		}
	}

	deprePerMonth := (car.Total - car.MinParfValue) / float64(car.Lifespan*utils.NumberOfMonthsInAYear)
	duration := car.Lifespan * utils.NumberOfMonthsInAYear

	var carEndDate time.Time
	if !carSoldDate.IsZero() {
		carEndDate = carSoldDate
	} else {
		carEndDate = carStartDate.AddDate(0, int(duration), 0)
	}

	assets := []db.Asset{}
	pv := car.Total
	for carStartDate.Before(carEndDate) {
		d := time.Date(carStartDate.Year(), carStartDate.Month(), 1, carStartDate.Hour(), 0, 0, 0, carStartDate.Location())
		d = utils.SetDateTo0000Hours(d)
		a := db.Asset{
			TransactionDate: utils.DateTime{Time: d},
			Type:            fmt.Sprintf(assetTypeCarFormat, car.Name),
			Amount:          pv,
		}

		assets = append(assets, a)
		carStartDate = carStartDate.AddDate(0, 1, 0)
		pv -= deprePerMonth
		pv = math.Max(pv, car.MinParfValue)
	}

	return assets, nil
}
