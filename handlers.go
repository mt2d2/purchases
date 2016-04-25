package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mt2d2/purchases/dal"
)

func writeJSON(w http.ResponseWriter, val interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(val); err != nil {
		log.Printf("error encoding json response: %s", err)
	}
}

func (app *app) handlePurchases(w http.ResponseWriter, req *http.Request) {
	purchases, err := dal.GetPurchasesAfterLastMarker(app.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, purchases)
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
	writeJSON(w, purchase)
}

func (app *app) handleDelete(w http.ResponseWriter, req *http.Request) {
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
	err = dal.DeletePurchase(app.db, byID)
	if err != nil {
		http.Error(w, "no such record", http.StatusNotFound)
		return
	}
}

func (app *app) handleAddPurchase(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var newPurchase dal.Purchase
	err := decoder.Decode(&newPurchase)
	if err != nil {
		log.Printf("could not decode json: %s", err)
		http.Error(w, "could not decode json", http.StatusInternalServerError)
		return
	}
	newPurchase.ID = 0
	newPurchase.TimeBought = time.Now()

	ok, errs := dal.ValidatePurchase(&newPurchase)
	if !ok || len(errs) != 0 {
		http.Error(w, "could not validate purchase", http.StatusInternalServerError)
		return
	}

	err = dal.AddPurchase(app.db, &newPurchase)
	if err != nil {
		http.Error(w, "could not save purchase", http.StatusInternalServerError)
		return
	}
}
