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
	carLoanDB db.CarLoanDB
	carConfig *carpb.CarConfig
}

func NewCarLoader(fw framework.FW) Loader {
	return &carLoader{
		fw:        fw,
		carLoanDB: fw.GetDB(db.CarLoansDatabaseName).(db.CarLoanDB),
		assetsDB:  fw.GetDB(db.AssetDatabaseName).(db.AssetDB),
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
	carLoans := []db.CarLoan{}
	for _, car := range l.carConfig.Cars {
		a, err := l.GetAssetsPerCar(car)
		if err != nil {
			return fmt.Errorf("could not get assets for car (%s): %s", car.Name, err)
		}

		cl, err := l.GetLiabilitiesPerCar(car)
		if err != nil {
			return fmt.Errorf("could not get liabilities for car (%s): %s", car.Name, err)
		}

		carLoans = append(carLoans, cl...)

		assetsMerged, err := l.mergeAssetWithLiabilitiesPerCar(a, cl)
		if err != nil {
			return fmt.Errorf("could merge assets with liability for car (%s): %s", car.Name, err)
		}

		carAssets = append(carAssets, assetsMerged...)
	}

	if err := l.assetsDB.BulkAdd(carAssets); err != nil {
		return fmt.Errorf("failed to bulk insert car assets into assets db: %s", err)
	}

	if err := l.carLoanDB.Clear(); err != nil {
		return fmt.Errorf("failed to clear car loan db: %s", err)
	}

	if err := l.carLoanDB.BulkAdd(carLoans); err != nil {
		return fmt.Errorf("failed to bulk insert car loans into car loan db: %s", err)
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

func (l *carLoader) mergeAssetWithLiabilitiesPerCar(assets []db.Asset, carLoans []db.CarLoan) ([]db.Asset, error) {
	carLoanPointer := 0
	mergedAssets := []db.Asset{}
	for _, a := range assets {
		if carLoanPointer == len(carLoans) {
			mergedAssets = append(mergedAssets, a)
			break
		}

		carLoan := carLoans[carLoanPointer]
		if carLoanMonthSameAsAsset(a, carLoan) {
			// merge
			a.Amount -= carLoan.AmountLeft
			carLoanPointer += 1
		} else {
			// Starting amount has the full loan amount unpaid
			a.Amount -= carLoans[0].AmountLeft + carLoans[0].AmountPaid
		}

		mergedAssets = append(mergedAssets, a)
	}
	return mergedAssets, nil
}

func carLoanMonthSameAsAsset(a db.Asset, carLoan db.CarLoan) bool {
	return carLoan.Date.Year() == a.TransactionDate.Year() && carLoan.Date.Month() == a.TransactionDate.Month()
}

// public for unit test
func (l *carLoader) GetLiabilitiesPerCar(car *carpb.Car) ([]db.CarLoan, error) {
	// Simple interest
	loan := car.Loan

	interest := loan.Amount * loan.InterestRate / 100.0 * float64(loan.Duration)
	totalPayable := loan.Amount + interest

	amountPaid := 0.0
	amountLeft := totalPayable

	paymentPerMonth := (totalPayable - loan.LastMonthAmount) / float64(int(loan.Duration)*12-1)

	date, err := utils.SetDateFromString(loan.StartDate)
	if err != nil {
		return nil, fmt.Errorf("unable to parse car loan start date (%s): %s", loan.StartDate, err)
	}

	loanSchedule := []db.CarLoan{}
	for i := 0; i < int(loan.Duration)*12-1; i++ {
		// Get start date
		amountPaid += paymentPerMonth
		amountLeft -= paymentPerMonth
		loan := db.CarLoan{
			Date:       date,
			Name:       car.Name,
			AmountPaid: amountPaid,
			AmountLeft: amountLeft,
		}
		loanSchedule = append(loanSchedule, loan)
		date = date.AddDate(0, 1, 0)
	}

	loanSchedule = append(loanSchedule, db.CarLoan{
		Date:       date,
		Name:       car.Name,
		AmountPaid: totalPayable,
		AmountLeft: 0.0,
	})

	return loanSchedule, nil
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

	// If car is bought after asset date, 1st of every month, skip
	first := carStartDate.Day() != 1

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

		carStartDate = carStartDate.AddDate(0, 1, 0)
		pv -= deprePerMonth
		pv = math.Max(pv, car.MinParfValue)
		if first {
			first = false
			continue
		}
		assets = append(assets, a)
	}

	return assets, nil
}
