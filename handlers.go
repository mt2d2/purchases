package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mt2d2/purchases/dal"
)

func (app *app) handlePurchases(w http.ResponseWriter, req *http.Request) {
	purchases, err := dal.GetPurchases(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(purchases)
}

func (app *app) handlePurchase(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "id feld is missing", http.StatusInternalServerError)
		return

	}
	byID, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	purchase, err := dal.GetPurchase(app.db, byID)
	if err != nil {
		http.Error(w, "no such record", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(purchase)
}
