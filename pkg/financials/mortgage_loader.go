package financials

import (
	"fmt"
	"math"

	"github.com/ashwinath/moneybags/pbgo/mortgagepb"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
)

const (
	numberOfMonthsInYear = 12
)

type mortgageLoader struct {
	fw                 framework.FW
	mortgageBulkLoader db.ClearAndBulkAdder
	mortgageConfig     *mortgagepb.MortgageConfig
}

func NewMortgageLoader(fw framework.FW) Loader {
	return &mortgageLoader{
		fw:                 fw,
		mortgageBulkLoader: fw.GetDB(db.MortgageDatabaseName).(db.ClearAndBulkAdder),
	}
}

func (l *mortgageLoader) Load() error {
	if err := l.loadMortgageConfig(); err != nil {
		return err
	}
	if err := l.mortgageBulkLoader.Clear(); err != nil {
		return fmt.Errorf("failed to clear mortgage db: %s", err)
	}

	for _, m := range l.mortgageConfig.Mortgages {
		if err := l.loadOneMortgageSchedule(m); err != nil {
			return fmt.Errorf("failed to load mortage schedule: %s", err)
		}
	}

	return nil
}

func (l *mortgageLoader) loadMortgageConfig() error {
	mortgageConfig := mortgagepb.MortgageConfig{}
	if err := utils.UnmarshalYAML(l.fw.GetConfig().FinancialsData.MortgageYamlFilepath, &mortgageConfig); err != nil {
		return fmt.Errorf("failed to unmarshal mortgage config: %s", err)
	}

	l.mortgageConfig = &mortgageConfig
	return nil
}

func (l *mortgageLoader) loadOneMortgageSchedule(m *mortgagepb.Mortgage) error {
	principal := m.Total
	for _, dp := range m.Downpayments {
		principal -= dp.Sum
	}

	monthlyPayment := CalculateMortgageMonthlyPayment(principal, m.InterestRatePercentage, int(m.MortgageDurationInYears))
	interestPaidSchedule := CalculateInterestPaidSchedule(principal, monthlyPayment, m.InterestRatePercentage)

	totalInterestToBePaid := 0.0
	for _, i := range interestPaidSchedule {
		totalInterestToBePaid += i
	}

	totalInterestLeft := totalInterestToBePaid
	totalPrincipalPaid := 0.0
	totalInterestPaid := 0.0
	totalPrincipalLeft := m.Total

	mortgageSchedule := []db.Mortgage{}

	// downpayment
	for _, dp := range m.Downpayments {
		totalPrincipalLeft -= dp.Sum
		totalPrincipalPaid += dp.Sum
		date, err := utils.SetDateFromString(dp.Date)
		if err != nil {
			return fmt.Errorf("could not parse downpayment date (%s): %s", dp.Date, err)
		}
		schedule := db.Mortgage{
			Date:               date,
			PrincipalPaid:      dp.Sum,
			TotalPrincipalPaid: totalPrincipalPaid,
			TotalInterestPaid:  totalInterestPaid,
			TotalPrincipalLeft: totalPrincipalLeft,
			TotalInterestLeft:  totalInterestLeft,
		}
		mortgageSchedule = append(mortgageSchedule, schedule)
	}

	mortgageDate, err := utils.SetDateFromString(m.MortgageFirstPayment)
	if err != nil {
		return fmt.Errorf("could not parse mortgage first payment date (%s): %s", m.MortgageFirstPayment, err)
	}

	for _, interestPaid := range interestPaidSchedule {
		// Interest
		totalInterestPaid += math.Min(totalInterestPaid+interestPaid, totalInterestToBePaid)
		totalInterestLeft = math.Max(totalInterestLeft-interestPaid, 0.0)

		// principal
		principalPaid := monthlyPayment - interestPaid
		totalPrincipalPaid = math.Min(totalPrincipalPaid+principalPaid, m.Total)
		totalPrincipalLeft = math.Max(totalPrincipalLeft-principalPaid, 0.0)

		schedule := db.Mortgage{
			Date:               mortgageDate,
			InterestPaid:       interestPaid,
			PrincipalPaid:      principalPaid,
			TotalPrincipalPaid: totalPrincipalPaid,
			TotalInterestPaid:  totalInterestPaid,
			TotalPrincipalLeft: totalPrincipalLeft,
			TotalInterestLeft:  totalInterestLeft,
		}
		mortgageSchedule = append(mortgageSchedule, schedule)
		mortgageDate = mortgageDate.AddDate(0, 1, 0)
	}

	if err := l.mortgageBulkLoader.BulkAdd(mortgageSchedule); err != nil {
		return fmt.Errorf("failed to bulk add mortgage schedule: %s", err)
	}

	return nil
}

// Public only because used in test
func CalculateMortgageMonthlyPayment(principal, interestRate float64, duration int) float64 {
	numberOfMonths := duration * numberOfMonthsInYear
	ir := interestRate / 100.0 / numberOfMonthsInYear
	// M = P [ i(1 + i)^n ] / [ (1 + i)^n – 1].
	return principal * (ir * (math.Pow(1.0+ir, float64(numberOfMonths)))) / (math.Pow(1.0+ir, float64(numberOfMonths)) - 1.0)
}

// Public only because used in test
func CalculateInterestPaidSchedule(principal, monthlyPayment, interestRatePercentage float64) []float64 {
	ir := interestRatePercentage / 100.0 / numberOfMonthsInYear
	sumLeft := principal

	interestPaidSchedule := []float64{}
	for sumLeft > 0.0 {
		interestPaid := sumLeft * ir
		interestPaidSchedule = append(interestPaidSchedule, interestPaid)
		sumLeft += interestPaid
		sumLeft -= monthlyPayment
	}

	return interestPaidSchedule
}
