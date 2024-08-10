package db

import (
	"time"

	"gorm.io/gorm"
)

const (
	SymbolDatabaseName string = "symbol"
	SymbolTypeStock    string = "stock"
	SymbolTypeCurrency string = "currency"
)

type Symbol struct {
	ID                uint `gorm:"primaryKey"`
	SymbolType        string
	Symbol            string
	BaseCurrency      *string
	LastProcessedDate *time.Time
}

type SymbolDB interface {
	CheckIfSymbolExists(symbol string) (bool, error)
	Insert(symbol *Symbol) error
	GetDistinctCurrencies() ([]string, error)
}

type symbolDB struct {
	db *gorm.DB
}

func NewSymbolDB(db *DB) (SymbolDB, error) {
	if err := db.DB.AutoMigrate(&Symbol{}); err != nil {
		return nil, err
	}

	return &symbolDB{
		db: db.DB,
	}, nil
}

func (db *symbolDB) CheckIfSymbolExists(symbol string) (bool, error) {
	var count int64
	res := db.db.Model(Symbol{}).Where("symbol = ?", symbol).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count == 1, nil
}

func (db *symbolDB) Insert(symbol *Symbol) error {
	result := db.db.Create(symbol)
	return result.Error
}

func (db *symbolDB) GetDistinctCurrencies() ([]string, error) {
	currencies := []string{}
	res := db.db.Distinct("base_currency").Model(Symbol{}).Where("symbol_type = ?", SymbolTypeStock).Select("base_currency").Find(&currencies)
	if res.Error != nil {
		return nil, res.Error
	}

	return currencies, nil
}
