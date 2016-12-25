package main

import (
	"boltapi/boltdb"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var database *boltdb.Database

func main() {
	dbfile := "/home/ubuntu/test.db"
	database, _ = boltdb.NewDatabase(dbfile)

	database.Put([]byte("bucket"), []byte("somekey"), []byte("somevalue"))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/read/", read)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func read(w http.ResponseWriter, r *http.Request) {
	log.Println(database)
	res, err := database.Get([]byte("bucket"), []byte("somekey"))
	log.Println("Result: ", string(res), "Err:", err)
}
