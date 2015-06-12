package main

import (
	"encoding/json"
	"net/http"

	"github.com/mt2d2/purchases/dal"
)

func (app *app) handlePurchases(w http.ResponseWriter, req *http.Request) {
	purchases, err := dal.GetPurchases(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(purchases)
}
