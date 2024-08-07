package db

import "time"

type Mortgage struct {
	ID                 uint      `gorm:"primaryKey"`
	Date               time.Time `gorm:"type:timestamp;unique"`
	InterestPaid       float64
	PrincipalPaid      float64
	TotalInterestPaid  float64
	TotalPrincipalPaid float64
	TotalInterestLeft  float64
	TotalPrincipalLeft float64
}

func (Mortgage) TableName() string {
	return "mortgage"
}
