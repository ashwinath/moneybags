package financials

import (
	"math"
	"testing"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
)

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
		assert.Equal(t, d, first.Date)
		assert.Equal(t, 0.0, first.InterestPaid)
		assert.Equal(t, 0.0, first.TotalInterestPaid)
		assert.Equal(t, 1000.0, first.PrincipalPaid)
		assert.Equal(t, 1000.0, first.TotalPrincipalPaid)
		assert.Equal(t, 49000.0, first.TotalPrincipalLeft)
		assert.True(t, math.Abs(10469.25-first.TotalInterestLeft) <= 0.01)

		firstPayment := mortgages[2]
		d, err = utils.SetDateFromString("2022-10-10")
		assert.Nil(t, err)
		assert.Equal(t, d, firstPayment.Date)
		assert.True(t, math.Abs(62.83-firstPayment.InterestPaid) <= 0.01)
		assert.True(t, math.Abs(62.83-firstPayment.TotalInterestPaid) <= 0.01)
		assert.True(t, math.Abs(68.73-firstPayment.PrincipalPaid) <= 0.01)
		assert.True(t, math.Abs(21068.73-firstPayment.TotalPrincipalPaid) <= 0.01)
		assert.True(t, math.Abs(28931.27-firstPayment.TotalPrincipalLeft) <= 0.01)
		assert.True(t, math.Abs(10406.41-firstPayment.TotalInterestLeft) <= 0.01)
	})

	assert.Nil(t, err)
}

// TODO: import test_mortgage over.
