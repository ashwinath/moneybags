package financials

import (
	"fmt"
	"testing"

	"github.com/ashwinath/moneybags/pbgo/carpb"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetAssetsPerCar(t *testing.T) {
	var tests = []struct {
		name             string
		car              *carpb.Car
		numberOfElements int
		deprePerMonth    float64
		lastValue        float64
		lastDate         string
	}{
		{
			name: "nominal, depre 10k/month, don't sell halfway",
			car: &carpb.Car{
				Name:         "Toy Car",
				Total:        140_000,
				MinParfValue: 20_000,
				Lifespan:     10,
				CarStartDate: "2020-01-01",
			},
			numberOfElements: 120,
			deprePerMonth:    1000.0,
			lastValue:        21_000,
			lastDate:         "2029-12-01",
		},
		{
			name: "nominal, depre 10k/month, sell halfway",
			car: &carpb.Car{
				Name:         "Taxi",
				Total:        140_000,
				MinParfValue: 20_000,
				Lifespan:     10,
				CarStartDate: "2020-01-01",
				CarSoldDate:  "2025-01-02",
			},
			numberOfElements: 61,
			deprePerMonth:    1000.0,
			lastValue:        80_000,
			lastDate:         "2025-01-01",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			l := carLoader{}
			assets, err := l.GetAssetsPerCar(tt.car)
			assert.Nil(t, err)
			assert.NotNil(t, assets)
			assert.Len(t, assets, tt.numberOfElements)

			// Check first value
			first := assets[0]
			assert.Equal(t, fmt.Sprintf("Car:%s", tt.car.Name), first.Type)
			assert.Equal(t, tt.car.Total, first.Amount)
			firstDate, err := utils.SetDateFromString(tt.car.CarStartDate)
			assert.Nil(t, err)
			assert.Equal(t, utils.DateTime{Time: firstDate}, first.TransactionDate)

			// Check second value
			second := assets[1]
			assert.Equal(t, fmt.Sprintf("Car:%s", tt.car.Name), second.Type)
			assert.Equal(t, tt.car.Total-tt.deprePerMonth, second.Amount)
			secondDate := firstDate.AddDate(0, 1, 0)
			assert.Equal(t, utils.DateTime{Time: secondDate}, second.TransactionDate)

			// Check last value
			last := assets[len(assets)-1]
			assert.Equal(t, fmt.Sprintf("Car:%s", tt.car.Name), last.Type)
			assert.Equal(t, tt.lastValue, last.Amount)
			lastDate, err := utils.SetDateFromString(tt.lastDate)
			assert.Nil(t, err)
			assert.Equal(t, utils.DateTime{Time: lastDate}, last.TransactionDate)
		})
	}
}
