package api

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type InventoryController struct{ db *ReadDb }

type AddInventoryDto struct{ Name string }

type CreateInventoryItem struct {
	ID   uuid.UUID
	Name string
}

var db *ReadDb
var bus *FakeBus

func NewInventoryController() InventoryController {
	db = NewReadDb()
	bus = NewFakeBus()
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

func (c InventoryController) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	item := AddInventoryDto{}

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	// bus.Send(CreateInventoryItem{uuid.NewV4(), item.Name})

}
