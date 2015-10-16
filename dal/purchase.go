package dal

import (
	"errors"
	"strings"
	"time"
)

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

// ValidatePurchase validates a purchase before database insertion
func ValidatePurchase(purchase *Purchase) (ok bool, errs []error) {
	errs = make([]error, 0)

	if strings.TrimSpace(purchase.Name) == "" {
		errs = append(errs, errors.New("Purchase must have a name."))
	}

	if purchase.Cost <= 0.0 {
		errs = append(errs, errors.New("Purchase must have a valid cost."))
	}

	if !purchase.TimeBought.Before(time.Now()) {
		errs = append(errs, errors.New("Purchase must have been made in the past."))
	}

	return len(errs) == 0, errs
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

// GetPurchasesAfterDate lists all purchases from the database made after date.
func GetPurchasesAfterDate(db *sqlx.DB, date time.Time) ([]Purchase, error) {
	purchases := []Purchase{}
	err := db.Select(&purchases,
		"SELECT * FROM purchase WHERE time_bought > ? ORDER BY time_bought DESC",
		date.Unix())
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

// DeletePurchase removes a purchase from the database.
func DeletePurchase(db *sqlx.DB, byID uint64) error {
	_, err := db.NamedExec(
		"DELETE FROM purchase where id=(:id)",
		map[string]interface{}{
			"id": byID,
		})
	return err
}
