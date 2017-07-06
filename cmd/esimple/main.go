package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jfcantin/esmadesimple/pkg/api"
)

const (
	port = "8000"
	addr = "127.0.0.1"
)

func main() {
	log.Println("Start Application")

	r := mux.NewRouter()

	ic := api.NewInventoryController()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/api/inventoryitems", ic.GetAllInventoryItems).Methods("GET")
	// r.HandleFunc("/api/inventoryitems", addInventoryItem).Methods("POST")

	// db = db.NewReadDb()

	srv := &http.Server{
		Handler:      r,
		Addr:         addr + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	log.Println("End Application")
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
