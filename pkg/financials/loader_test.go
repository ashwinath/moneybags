package financials

import (
	"path"
	"testing"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"github.com/ashwinath/moneybags/pkg/config"
	database "github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/framework"
	"github.com/ashwinath/moneybags/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func createFW(t *testing.T, db *database.DB) framework.FW {
	assetDB, err := database.NewAssetDB(db)
	assert.Nil(t, err)

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	assert.Nil(t, err)

	subsituteLocalRepoLocation(c)

	return framework.New(c, sugar, map[string]any{
		database.AssetDatabaseName: assetDB,
	})
}

func subsituteLocalRepoLocation(c *configpb.Config) {
	p := c.FinancialsData.AssetsCsvFilepath
	newPath := path.Join(utils.GetLocalRepoLocation(), p)
	c.FinancialsData.AssetsCsvFilepath = newPath
}

func TestLoad(t *testing.T) {
	err := database.RunTest(func(db *database.DB) {
		fw := createFW(t, db)
		loader := NewLoader(fw)
		err := loader.Start()
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}
