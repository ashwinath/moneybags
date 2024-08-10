package financials

import (
	"testing"

	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestStocksLoader(t *testing.T) {

	err := database.RunTest(func(db *database.DB) {
		_ = createFW(t, db)
		// TODO: How to test this?
		// mock interface of alphavantage
	})
	assert.Nil(t, err)
}
