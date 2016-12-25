package main

import (
	"boltapi/boltdb"
	"io/ioutil"
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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/dbs/{db}/buckets/{bucket}/keys/{key}", reqHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		put(w, r)
	case "GET":
		get(w, r)
	}
	vars := mux.Vars(r)
	log.Println("db:", vars["db"])
}

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]
	res, err := database.Get([]byte(bucket), []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	response := "{" + key + ":" + string(res) + "}"
	w.Write([]byte(response))
}

func put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]
	val, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	database.Put([]byte(bucket), []byte(key), []byte(val))
}
