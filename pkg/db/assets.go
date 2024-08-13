package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const AssetDatabaseName string = "asset"

type Asset struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamp;uniqueIndex:ix_date_type_assets" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_assets" csv:"type"`
	Amount          float64        `csv:"amount"`
}

type AssetDB interface {
	BulkAdd(objs interface{}) error
}

type assetDB struct {
	db *gorm.DB
}

func NewAssetDB(db *DB) (AssetDB, error) {
	if err := db.DB.AutoMigrate(&Asset{}); err != nil {
		return nil, err
	}

	return &assetDB{
		db: db.DB,
	}, nil
}

// Clears the database
func (db *assetDB) Clear() error {
	return db.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Asset{}).Error
}

// Bulk add data
func (db *assetDB) BulkAdd(objs interface{}) error {
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(objs).Error
}

func (db *assetDB) Count() (int64, error) {
	var count int64
	r := db.db.Model(&Asset{}).Count(&count)
	if r.Error != nil {
		return 0, r.Error
	}
	return count, nil
}
