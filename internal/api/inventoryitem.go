package api

import uuid "github.com/satori/go.uuid"

type InventoryItemCreated struct {
	ID   uuid.UUID
	Name string
}

// InventoryItem represents
type InventoryItem struct {
	ID      uuid.UUID
	Version int
	changes []interface{}
}

func (item *InventoryItem) LoadFromHistory(h []interface{}) {
	for e := range h {
		item.apply(e)
		item.Version++
	}
}

func (item *InventoryItem) apply(e interface{}) {
	switch t := e.(type) {
	case InventoryItemCreated:
		item.ID = t.ID
	}
}
