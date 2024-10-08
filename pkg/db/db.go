package db

import (
	"fmt"
	"time"

	"github.com/ashwinath/moneybags/pbgo/configpb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	DB *gorm.DB
}

// New initialises a new base database object.
func NewBaseDB(dbConfig *configpb.PostgresDB) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Singapore",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DbName,
		dbConfig.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation("Asia/Singapore")
			return time.Now().In(ti)
		},
	})

	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, nil
}

func (d *DB) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
