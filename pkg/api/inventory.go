package api

import (
	"encoding/json"
	"net/http"
)

type (
	InventoryController struct {
		db *ReadDb
	}
)

var db *ReadDb

func NewInventoryController() InventoryController {
	db = NewReadDb()
	return InventoryController{db: db}
}

func (c InventoryController) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items := c.db.GetInventoryItems()
	b, err := json.Marshal(items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// func addInventoryItem(w http.ResponseWriter, r *http.Request) {
// 	// params := mux.Vars(r)
// 	item := InventoryItem{}
// 	err := json.NewDecoder(r.Body).Decode(&item)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 	}
// }
