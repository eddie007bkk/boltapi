package main

import (
	"boltapi/boltdb"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
)

type userRequest struct {
	db, bucket, key string
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var database *boltdb.Database

func main() {
	port := ":8080"
	log.Println("Server Listening on Port: ", port)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/dbs/", requestHandler)
	router.HandleFunc("/dbs/{db}/", requestHandler)
	router.HandleFunc("/dbs/{db}/stats/", requestHandler)
	router.HandleFunc("/dbs/{db}/compact/", requestHandler)
	router.HandleFunc("/dbs/{db}/buckets/{bucketName}/keys/{keyName}", requestHandler)

	log.Fatal(http.ListenAndServe(port, router))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var userRequest userRequest
	userRequest.db = vars["db"]
	userRequest.bucket = vars["bucketName"]
	userRequest.key = vars["keyName"]

	if len(userRequest.db) > 0 && (database == nil || (userRequest.db != database.CurrentDB())) {
		var err error
		database, err = openDB(userRequest.db)
		handleErr(err)
	}

	switch r.Method {
	case "PUT":
		put(w, r, userRequest)
	case "GET":
		get(w, r, userRequest)
	case "DELETE":
		delete(w, r, userRequest)
	}
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

func getCurrentDB(w http.ResponseWriter, r *http.Request, userRequest userRequest) {
	if database == nil {
		log.Println("DB = Nil")
		w.Write([]byte("{\"database\":\"none\"}"))
	} else {
		dbPath := database.CurrentDB()
		w.Write([]byte("{\"database\":\"" + dbPath + "\"}"))
	}
}

func get(w http.ResponseWriter, r *http.Request, userRequest userRequest) {
	log.Println(userRequest)
	if len(userRequest.key) > 0 {
		//User has specified a db, bucket & key
		res, err := database.Get([]byte(userRequest.bucket), []byte(userRequest.key))
		handleErr(err)
		response := "{\"" + userRequest.key + "\":\"" + string(res) + "\"}"
		w.Write([]byte(response))
	} else if len(userRequest.db) > 0 {
		//User has only specified a db, return data about this db
		res := database.BK.Stats()
		log.Println(reflect.TypeOf(res), res)

	} else {
		//User has not specified any db, return data about all DBs
	}

}

func put(w http.ResponseWriter, r *http.Request, userRequest userRequest) {
	if len(userRequest.key) > 0 {
		//User has specified a db, bucket & key
		vars := mux.Vars(r)
		bucket := vars["bucket"]
		key := vars["key"]
		val, err := ioutil.ReadAll(r.Body)
		handleErr(err)
		database.Put([]byte(bucket), []byte(key), []byte(val))
	} else if len(userRequest.db) > 0 {
		//User has only specified a db, open a new database
	} else {
		//User has not specified any db, do nothing
	}

}
func delete(w http.ResponseWriter, r *http.Request, userRequest userRequest) error {
	if database != nil {
		err := database.DB.Close()
		handleErr(err)
	}
	dbFolder := "/home/ubuntu/"
	dbPath := dbFolder + userRequest.db
	err := os.Remove(dbPath)
	handleErr(err)
	return err
}
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
