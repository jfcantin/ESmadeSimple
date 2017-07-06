package api

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
)

type ReadDb struct {
	// list map[uuid.UUID]InventoryItemListDto
	list []InventoryItemListDto
}

func NewReadDb() *ReadDb {
	// list := make(map[uuid.UUID]InventoryItemListDto)
	var items []InventoryItemListDto
	for i := 0; i < 6; i++ {
		id := uuid.NewV4()
		items = append(items, InventoryItemListDto{id, "Item " + strconv.Itoa(i)})
	}
	return &ReadDb{list: items}
}

func (db *ReadDb) GetInventoryItems() []InventoryItemListDto {
	return db.list
}

type InventoryItemListDto struct {
	ID   uuid.UUID
	Name string
}
