package api

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestLoadingFromHistory(t *testing.T) {
	id := uuid.NewV4()
	history := []interface{}{InventoryItemCreated{ID: id, Name: "boo"}}

	var item InventoryItem
	item.LoadFromHistory(history)

	if item.Version != 1 {
		t.Errorf("Expected %v, but was %v", 1, item.Version)
	}
}
