package financials

import (
	"fmt"
	"math"
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pbgo/mortgagepb"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/ashwinath/simple/framework"
)

type mortgageLoader struct {
	fw                 framework.FW
	mortgageDB         db.MortgageDB
	mortgageBulkLoader db.ClearAndBulkAdder
	mortgageConfig     *mortgagepb.MortgageConfig
}

func NewMortgageLoader(fw framework.FW) Loader {
	return &mortgageLoader{
		fw:                 fw,
		mortgageBulkLoader: fw.GetDB(db.MortgageDatabaseName).(db.ClearAndBulkAdder),
		mortgageDB:         fw.GetDB(db.MortgageDatabaseName).(db.MortgageDB),
	}
}

func (mortgageLoader) Name() string {
	return "mortgage"
}

func (l *mortgageLoader) Load() error {
	if err := l.loadMortgageConfig(); err != nil {
		return err
	}
	if err := l.mortgageBulkLoader.Clear(); err != nil {
		return fmt.Errorf("failed to clear mortgage db: %s", err)
	}

	for _, group := range l.mortgageConfig.Groups {
		for _, m := range group.Mortgages {
			if err := l.loadOneMortgageSchedule(m, group); err != nil {
				return fmt.Errorf("failed to load mortage schedule: %s", err)
			}
		}
	}

	// TODO: Query all mortgages and calculate the interest paid and principal paid
	if err := l.calculateTotalPrincipalAndInterestPaid(); err != nil {
		return fmt.Errorf("failed to calculate mortgage total principal and interest: %s", err)
	}

	return nil
}

type cumulativeGroup struct {
	TotalPrincipalPaid float64
	TotalInterestPaid  float64
}

func (l *mortgageLoader) calculateTotalPrincipalAndInterestPaid() error {
	mortgages, err := l.mortgageDB.GetMortgage()
	if err != nil {
		return fmt.Errorf("unable to query mortgage: %s", err)
	}

	newMortgages := []db.Mortgage{}
	cumulative := map[string]*cumulativeGroup{}
	for _, m := range mortgages {
		if _, ok := cumulative[m.GroupName]; !ok {
			cumulative[m.GroupName] = &cumulativeGroup{}
		}
		cg := cumulative[m.GroupName]
		cg.TotalInterestPaid += m.InterestPaid
		cg.TotalPrincipalPaid += m.PrincipalPaid

		m.TotalInterestPaid = cg.TotalInterestPaid
		m.TotalPrincipalPaid = cg.TotalPrincipalPaid

		newMortgages = append(newMortgages, m)
	}

	if err := l.mortgageDB.BulkUpdate(newMortgages); err != nil {
		return fmt.Errorf("failed to bulk update mortgage schedule: %s", err)
	}

	return nil
}

func (l *mortgageLoader) loadMortgageConfig() error {
	mortgageConfig := mortgagepb.MortgageConfig{}
	if err := utils.UnmarshalYAML(l.fw.GetConfig().(*configpb.Config).FinancialsData.MortgageYamlFilepath, &mortgageConfig); err != nil {
		return fmt.Errorf("failed to unmarshal mortgage config: %s", err)
	}

	l.mortgageConfig = &mortgageConfig
	return nil
}

func (l *mortgageLoader) loadOneMortgageSchedule(m *mortgagepb.Mortgage, mg *mortgagepb.MortgageGroup) error {
	principal := m.Total
	for _, dp := range m.Downpayments {
		principal -= dp.Sum
	}

	monthlyPayment := CalculateMortgageMonthlyPayment(principal, m.InterestRatePercentage, int(m.MortgageDurationInYears))
	interestPaidSchedule := CalculateInterestPaidSchedule(principal, monthlyPayment, m.InterestRatePercentage)

	mortgageSchedule := []db.Mortgage{}
	// downpayment
	for _, dp := range m.Downpayments {
		date, err := utils.SetDateFromString(dp.Date)
		if err != nil {
			return fmt.Errorf("could not parse downpayment date (%s): %s", dp.Date, err)
		}
		schedule := db.Mortgage{
			Date:          date,
			PrincipalPaid: dp.Sum,
			GroupName:     mg.Name,
		}
		mortgageSchedule = append(mortgageSchedule, schedule)
	}

	mortgageDate, err := utils.SetDateFromString(m.MortgageFirstPayment)
	if err != nil {
		return fmt.Errorf("could not parse mortgage first payment date (%s): %s", m.MortgageFirstPayment, err)
	}

	hasEndDate := false
	var endDate time.Time
	if m.MortgageEndDate != nil {
		hasEndDate = true
		var err error
		endDate, err = utils.SetDateFromString(*m.MortgageEndDate)
		if err != nil {
			return fmt.Errorf("could not parse mortgage end date (%s): %s", *m.MortgageEndDate, err)
		}
	}

	for _, interestPaid := range interestPaidSchedule {
		if hasEndDate {
			ed := utils.SetDateToEndOfMonth(endDate)
			if mortgageDate.After(ed) {
				break
			}
		}

		schedule := db.Mortgage{
			Date:          mortgageDate,
			InterestPaid:  interestPaid,
			PrincipalPaid: monthlyPayment - interestPaid,
			GroupName:     mg.Name,
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
	numberOfMonths := duration * utils.NumberOfMonthsInAYear
	ir := interestRate / 100.0 / utils.NumberOfMonthsInAYear
	// M = P [ i(1 + i)^n ] / [ (1 + i)^n – 1].
	return principal * (ir * (math.Pow(1.0+ir, float64(numberOfMonths)))) / (math.Pow(1.0+ir, float64(numberOfMonths)) - 1.0)
}

// Public only because used in test
func CalculateInterestPaidSchedule(principal, monthlyPayment, interestRatePercentage float64) []float64 {
	ir := interestRatePercentage / 100.0 / utils.NumberOfMonthsInAYear
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
