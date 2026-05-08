package financials

import (
	"fmt"
	"testing"
	"time"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var singaporeLocation, _ = time.LoadLocation("Asia/Singapore")

func TestMortgageSchedule(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)
		ml := NewMortgageLoader(fw)
		err := ml.Load()
		assert.Nil(t, err)

		mortgages := []database.Mortgage{}
		res := db.DB.Order("date asc").Find(&mortgages)
		assert.Nil(t, res.Error)

		assert.Greater(t, len(mortgages), 1)

		first := mortgages[0]
		d, err := utils.SetDateFromString("2021-10-10")
		assert.Nil(t, err)
		assert.Equal(t, d.In(singaporeLocation), first.Date.In(singaporeLocation))
		assert.Equal(t, 0.0, first.InterestPaid)
		assert.Equal(t, 0.0, first.TotalInterestPaid)
		assert.Equal(t, 1000.0, first.PrincipalPaid)
		assert.Equal(t, 1000.0, first.TotalPrincipalPaid)

		firstPayment := mortgages[2]
		d, err = utils.SetDateFromString("2022-10-10")
		assert.Nil(t, err)
		assert.Equal(t, d.In(singaporeLocation), firstPayment.Date.In(singaporeLocation))
		assert.InDelta(t, 62.83, firstPayment.InterestPaid, 0.1)
		assert.InDelta(t, 62.83, firstPayment.TotalInterestPaid, 0.1)
		assert.InDelta(t, 68.73, firstPayment.PrincipalPaid, 0.1)
		assert.InDelta(t, 21068.73, firstPayment.TotalPrincipalPaid, 0.1)

		// Test house asset loader also
		hal := NewHouseAssetLoader(fw)
		err = hal.Load()
		assert.Nil(t, err)

		var count int64
		res = db.DB.Model(database.Asset{}).
			Where("type = ?", "House").
			Count(&count)
		assert.Nil(t, res.Error)
		assert.Greater(t, count, int64(1))
	})

	assert.Nil(t, err)
}

func TestMortgageAndInterest(t *testing.T) {
	var tests = []struct {
		principal              float64
		ir                     float64
		years                  int
		expectedMonthlyPayment float64
		expectedInterestPaid   float64
	}{
		{
			principal:              29_000.0,
			ir:                     2.6,
			years:                  25,
			expectedMonthlyPayment: 131.56,
			expectedInterestPaid:   10469.25,
		},
		{
			principal:              500_000.0,
			ir:                     5.0,
			years:                  35,
			expectedMonthlyPayment: 2523.44,
			expectedInterestPaid:   559844.12,
		},
		{
			principal:              25_321_323.0,
			ir:                     1.2,
			years:                  20,
			expectedMonthlyPayment: 118_724.61,
			expectedInterestPaid:   3_172_582.59,
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			monthlyPayment := CalculateMortgageMonthlyPayment(
				tt.principal, tt.ir, tt.years,
			)
			assert.InDelta(t, monthlyPayment, tt.expectedMonthlyPayment, 0.1)

			interestPaidSchedule := CalculateInterestPaidSchedule(
				tt.principal, monthlyPayment, tt.ir,
			)

			sum := 0.0
			for _, ip := range interestPaidSchedule {
				sum += ip
			}

			assert.InDelta(t, sum, tt.expectedInterestPaid, 0.1)
		})
	}
}
