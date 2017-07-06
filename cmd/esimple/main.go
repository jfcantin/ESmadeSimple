package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const (
	port = "8000"
	addr = "127.0.0.1"
)

var db *ReadDb

type ReadDb struct {
	// list map[uuid.UUID]InventoryItemListDto
	list []InventoryItemListDto
}

func (db *ReadDb) NewReadDb() *ReadDb {
	// list := make(map[uuid.UUID]InventoryItemListDto)
	var items []InventoryItemListDto
	for i := 0; i < 5; i++ {
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

func main() {
	log.Println("Start Application")
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/inventoryitems", getAllInventoryItems).Methods("GET")
	// r.HandleFunc("/api/inventoryitems", addInventoryItem).Methods("POST")

	db = db.NewReadDb()

	srv := &http.Server{
		Handler:      r,
		Addr:         addr + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	log.Println("End Application")
}

func getAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	items := db.GetInventoryItems()
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home\n"))
}
