package dal

import "time"

import (
	"github.com/jmoiron/sqlx"
)

// Purchase is the primary dal for an item purchased.
type Purchase struct {
	ID         uint64    `json:"id"`
	Name       string    `json:"name"`
	Cost       float64   `json:"cost"`
	TimeBought time.Time `json:"time_bought" db:"time_bought"`
}

// AddPurchase adds a new a purchase to the database.
func AddPurchase(db *sqlx.DB, purchase *Purchase) error {
	_, err := db.NamedExec(
		"INSERT INTO purchase (name, cost, time_bought) VALUES (:name, :cost, :time_bought)",
		map[string]interface{}{
			"name":        purchase.Name,
			"cost":        purchase.Cost,
			"time_bought": purchase.TimeBought,
		})
	return err
}

// GetPurchases lists all purchases from the database.
func GetPurchases(db *sqlx.DB) ([]Purchase, error) {
	purchases := []Purchase{}
	err := db.Select(&purchases, "SELECT * FROM purchase ORDER BY time_bought DESC")
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

// GetPurchase retrieves a purchase by its ID.
func GetPurchase(db *sqlx.DB, byID uint64) (*Purchase, error) {
	purchase := Purchase{}
	err := db.Get(&purchase, "SELECT * FROM purchase WHERE id=?", byID)
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}
