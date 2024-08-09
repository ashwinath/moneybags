package db

import (
	"time"

	"gorm.io/gorm"
)

const AssetDatabaseName string = "asset"

type Asset struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionDate time.Time `gorm:"type:timestamp;uniqueIndex:ix_date_type_assets"`
	Type            string    `gorm:"uniqueIndex:ix_date_type_assets"`
	Amount          float64
}

type AssetDB interface {
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
