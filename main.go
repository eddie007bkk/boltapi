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

func main() {
	dbfile := "/home/ubuntu/test.db"
	database, err := boltdb.NewDatabase(dbfile)
	if err != nil {
		log.Println(err)
	}

	database.Put([]byte("bucket"), []byte("somekey"), []byte("somevalue"))
	res, err := database.Get([]byte("bucket"), []byte("somekey"))
	log.Println("Result: ", string(res), "Err:", err)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Hello)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello User"))
}
