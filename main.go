package main

import (
	"boltapi/db"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)



func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var database *db.Database

const dbsFolder string = "/home/ubuntu/boltdbs"

func main() {

	if _, err := os.Stat(dbsFolder); os.IsNotExist(err) {
		os.Mkdir(dbsFolder, os.ModePerm)
	}

	port := ":8080"
	log.Println("Server Listening on Port: ", port)

	router := mux.NewRouter().StrictSlash(true)

	//txHandler functions are transactional requests to insert, update or delete entries
	router.HandleFunc("/dbs/{db}/", txHandler)
	router.HandleFunc("/dbs/{db}/buckets/{bucketName}", txHandler)
	router.HandleFunc("/dbs/{db}/buckets/{bucketName}/keys/{keyName}", txHandler)

	//adminHandler are requests for information about the database or an action such as compaction
	router.HandleFunc("/dbs/", adminHandler)
	router.HandleFunc("/dbs/{db}/stats/", adminHandler)
	router.HandleFunc("/dbs/{db}/compact/", adminHandler)
	router.HandleFunc("/dbs/{db}/buckets/", adminHandler)
	//router.HandleFunc("/dbs/{db}/buckets/{bucketName}/keys", requestHandler)

	log.Fatal(http.ListenAndServe(port, router))
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
