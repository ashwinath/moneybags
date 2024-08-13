package db

import (
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/utils"
	"go.uber.org/zap"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	r           = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func createTestDatabase() (*string, error) {
	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	if err != nil {
		return nil, err
	}

	db, err := NewBaseDB(c.PostgresDb)
	if err != nil {
		return nil, err
	}

	databaseName := fmt.Sprintf("%s_%s", "test", randStringRunes(10))

	// Drop database if exists
	stmt := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", databaseName)
	if rs := db.DB.Exec(stmt); rs.Error != nil {
		return nil, rs.Error
	}

	// if not create it
	stmt = fmt.Sprintf("CREATE DATABASE %s;", databaseName)
	if rs := db.DB.Exec(stmt); rs.Error != nil {
		return nil, rs.Error
	}

	// close db connection
	sql, err := db.DB.DB()
	defer func() {
		_ = sql.Close()
	}()
	if err != nil {
		return nil, err
	}

	return &databaseName, nil
}

func dropDatabase(name string) error {
	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	if err != nil {
		return err
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := NewBaseDB(c.PostgresDb)
	if err != nil {
		return err
	}

	// Drop database if exists
	stmt := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", name)
	if rs := db.DB.Exec(stmt); rs.Error != nil {
		return rs.Error
	}

	return nil
}

func RunTest(function func(db *DB)) error {
	dbName, err := createTestDatabase()
	if err != nil {
		return err
	}

	p := path.Join(utils.GetLocalRepoLocation(), "./pkg/config/testdata/config.yaml")
	c, err := config.New(p)
	if err != nil {
		return err
	}

	c.PostgresDb.DbName = *dbName

	db, err := NewBaseDB(c.PostgresDb)
	if err != nil {
		return err
	}

	defer dropDatabase(*dbName)
	defer db.Close()

	function(db)

	return nil
}
