package dal

import "time"

import (
	"github.com/jmoiron/sqlx"
)

// Purchase is the primary dal for an item purchased.
type Purchase struct {
	ID         uint64
	Name       string
	Cost       float64
	TimeBought time.Time `db:"time_bought"`
}

// NewPurchase adds a new a purchase to the database.
func NewPurchase(db *sqlx.DB) error {
	return nil
}

// GetPurchases lists all purchases from the database.
func GetPurchases(db *sqlx.DB) ([]Purchase, error) {
	purchases := []Purchase{}
	err := db.Select(&purchases, "SELECT * FROM purchase ORDER BY time_bought ASC")
	if err != nil {
		return nil, err
	}
	return purchases, nil
}
