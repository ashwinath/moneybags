package db

import (
	"github.com/ashwinath/moneybags/pkg/utils"
	"gorm.io/gorm"
)

const AssetDatabaseName string = "asset"

type Asset struct {
	ID              uint           `gorm:"primaryKey"`
	TransactionDate utils.DateTime `gorm:"type:timestamp;uniqueIndex:ix_date_type_assets" csv:"date"`
	Type            string         `gorm:"uniqueIndex:ix_date_type_assets" csv:"type"`
	Amount          float64        `csv:"amount"`
}

type AssetDB interface {
	Clear() error
	BulkAdd(assets []*Asset) error
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

func (db *assetDB) BulkAdd(assets []*Asset) error {
	return db.db.Create(assets).Error
}
