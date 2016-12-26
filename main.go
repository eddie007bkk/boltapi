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

	log.Println("Begin Server")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/dbs/{db}/buckets/{bucket}/keys/{key}", reqHandler)
	router.HandleFunc("/dbs/current", getCurrentDB)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func openDB(dbfile string) (d *boltdb.Database, err error) {
	if database != nil {
		err := database.DB.Close()
		handleErr(err)
	}
	dbFolder := "/home/ubuntu/"
	dbPath := dbFolder + dbfile
	database, err := boltdb.NewDatabase(dbPath)
	handleErr(err)
	return database, err
}

func getCurrentDB(w http.ResponseWriter, r *http.Request) {
	if database == nil {
		log.Println("DB = Nil")
		w.Write([]byte("{\"database\":\"none\"}"))
	} else {
		dbPath := database.CurrentDB()
		w.Write([]byte("{\"database\":\"" + dbPath + "\"}"))
	}
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestDB := vars["db"]

	if database == nil || (requestDB != database.CurrentDB()) {
		var err error
		database, err = openDB(requestDB)
		handleErr(err)
	}
	switch r.Method {
	case "PUT":
		put(w, r)
	case "GET":
		get(w, r)
	}
	//vars := mux.Vars(r)
	//log.Println("db:", vars["db"])
}

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]
	res, err := database.Get([]byte(bucket), []byte(key))
	handleErr(err)
	response := "{\"" + key + "\":\"" + string(res) + "\"}"
	w.Write([]byte(response))
}

func put(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]
	val, err := ioutil.ReadAll(r.Body)
	handleErr(err)
	database.Put([]byte(bucket), []byte(key), []byte(val))
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
