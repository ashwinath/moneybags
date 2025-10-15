package financials

import (
	"fmt"
	"math"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/simple/framework"
)

type incomeLoader struct {
	fw                  framework.FW
	incomeDB            db.IncomeDB
	totalCompensationDB db.TotalCompensationDB
}

func NewIncomeLoader(fw framework.FW) Loader {
	return &incomeLoader{
		fw:                  fw,
		incomeDB:            fw.GetDB(db.IncomeDatabaseName).(db.IncomeDB),
		totalCompensationDB: fw.GetDB(db.TotalCompensationName).(db.TotalCompensationDB),
	}
}

func (incomeLoader) Name() string {
	return "income"
}

func (l *incomeLoader) Load() error {
	clearAndBulkAdder := l.totalCompensationDB.(db.ClearAndBulkAdder)
	if err := clearAndBulkAdder.Clear(); err != nil {
		return fmt.Errorf("unable to clear total_compensation db: %s", err)
	}

	/* need to create a table with the following
	   year, gross base, gross bonus, rsu, espp, gross tc, cpf, tc
	   income types:
	   - Base
	   - Base Bonus
	   - Overtime
	   - CPF
	   - CPF Overtime
	   - CPF Bonus
	   - ESPP
	   - RSU
	*/
	types := map[string][]string{
		"base":  {"Base"},
		"bonus": {"Base Bonus", "Overtime"},
		"cpf":   {"CPF", "CPF Overtime", "CPF Bonus"},
		"espp":  {"ESPP"},
		"rsu":   {"RSU"},
	}
	results := map[string][]db.AggregateIncomeByYear{}

	for t, dbTypes := range types {
		rows, err := l.incomeDB.AggregateByYear(dbTypes)
		if err != nil {
			return fmt.Errorf("unable to query db for %s type, %s", t, err)
		}
		results[t] = rows
	}

	// Since I started working from 2016, will initialise from 2016
	tcs := map[int]*db.TotalCompensation{}

	for i := 2016; i <= time.Now().Year(); i++ {
		tc := db.TotalCompensation{
			Year: i,
		}
		tcs[i] = &tc
	}

	// transform
	for _, agg := range results["base"] {
		cpf := getCPFAmountFromBase(agg.Year, agg.Amount)
		tc := tcs[agg.Year]
		tc.GrossBase += agg.Amount + cpf

		// CPF is out, easier to calc
		tc.TotalCompensation += agg.Amount
	}

	for _, agg := range results["bonus"] {
		cpf := agg.Amount / 80 * 20
		tc := tcs[agg.Year]
		tc.GrossBonus = agg.Amount + cpf

		// CPF is out, easier to calc
		tc.TotalCompensation += agg.Amount
	}

	for _, agg := range results["cpf"] {
		tc := tcs[agg.Year]
		tc.CPF = agg.Amount
	}

	for _, agg := range results["espp"] {
		tc := tcs[agg.Year]
		tc.ESPP = agg.Amount
	}

	for _, agg := range results["rsu"] {
		tc := tcs[agg.Year]
		tc.RSU = agg.Amount
	}

	for _, tc := range tcs {
		tc.GrossTC = tc.GrossBase + tc.GrossBonus + tc.RSU + tc.ESPP
		tc.TotalCompensation += tc.RSU + tc.ESPP + tc.CPF
	}

	if err := clearAndBulkAdder.BulkAdd(mapToSlice(tcs)); err != nil {
		return fmt.Errorf("unable to bulk add total compensations: %s", err)
	}

	return nil
}

func getCPFAmountFromBase(year int, amount float64) float64 {
	// we will just ignore the sep23 - jan 24 6300 rate, makes life too hard
	// cpf is 20%, so amount here is 80%
	cpfUncapped := amount / 80 * 20
	var cpf float64
	if year < 2024 {
		cpf = math.Min(6800.00*0.2, cpfUncapped)
	} else if year == 2025 {
		cpf = math.Min(7400.00*0.2, cpfUncapped)
	} else {
		cpf = math.Min(8000.00*0.2, cpfUncapped)
	}

	return cpf
}

func mapToSlice(tcs map[int]*db.TotalCompensation) []*db.TotalCompensation {
	var slice []*db.TotalCompensation
	for _, tc := range tcs {
		slice = append(slice, tc)
	}
	return slice
}
