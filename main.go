package main

import (
	"boltapi/db"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

type userRequest struct {
	db, bucket, key, cmd string
}

func (userReq *userRequest) GetRequest(r *http.Request) {
	vars := mux.Vars(r)
	userReq.db = vars["db"]
	userReq.bucket = vars["bucketName"]
	userReq.key = vars["keyName"]
	return
}

func (userReq *userRequest) GetUserDB(r *http.Request) {
	if len(userReq.db) > 0 && (database == nil || (userReq.db != database.CurrentDB())) {
		var err error
		database, err = openDB(userReq.db)
		handleErr(err)
	}
	//Otherwise, system is already using user db.
}

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

	//TODO: Router paths need to be cleaned up with regular expressions/options
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

func txHandler(w http.ResponseWriter, r *http.Request) {
	var userReq userRequest
	userReq.GetRequest(r)
	log.Println(userReq.db, userReq.bucket, userReq.key, userReq.key)
	if len(userReq.key) > 0 {
		switch r.Method {
		case "PUT":
			put(w, r, userReq)
		case "GET":
			get(w, r, userReq)
		case "DELETE":
			delete(w, r, userReq)
		}
		/*
			If we have a keyName, then the user has also specified a dbs and bucket
			Possible Actions:
				GET - Read a value given a key in the URL
				PUT - Insert a value from the body, given a key in the URL
				DELETE - Delete a key/value pair, given a key in the URL
		*/
	} else if len(userReq.bucket) > 0 {
		/*
			We have a bucketName but not a keyName
			Possible Actions:
				DELETE - Delete bucket & all contents.
		*/
	} else if len(userReq.db) > 0 {
		/*
			We only have a database name.
			Possible actions:
				DELETE - Delite entire database
		*/
	}
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	//Get variables {db}, {bucketName}, and {keyName} from user request URL
	var userRequest userRequest
	userRequest.GetRequest(r)

	reqURI := r.URL.RequestURI()
	reqURI = reqURI[1 : len(reqURI)-1]
	uri := strings.Split(reqURI, "/")

	if (len(uri) % 2) == 0 {
		//Even number of entries in uri means we ended with either a specific db, bucket or key
		log.Println("Even URI")

	} else {
		/*Odd number of entries in uri means we ended with either dbs, buckets or keys (general)
		Possible Commands:
			/dbs/
				- GET - Show All Databases
			/dbs/{database}/stats
				- GET - Return stats grid about database
			/dbs/{database}/compact
				- POST - Compact database by reading & rewriting entire dbs
			/dbs/{database}/buckets
				- GET - Show All Buckets
			/dbs/{database}/buckets/{bucket}/keys/
				- GET - Show all keys in bucket
		*/
		userRequest.cmd = uri[len(uri)-1]
		switch userRequest.cmd {
		case "dbs":
			log.Println("Show all databases")
		case "stats":
			log.Println("Show stats for database: ", userRequest.db)
		case "compact":
			log.Println("Compact this database:", userRequest.db)
		case "buckets":
			log.Println("Show all buckets in this database: ", userRequest.db)
		case "keys":
			log.Println("Show all keys in database:", userRequest.db, "bucket:", userRequest.bucket)
		default:
			log.Println("Not possible because routes are predefined by router.")
		}

	}

	w.Write([]byte("Bananabananabanana"))

}

func openDB(dbfile string) (d *db.Database, err error) {
	if database != nil {
		err := database.DB.Close()
		handleErr(err)
	}
	dbPath := dbsFolder + "/" + dbfile
	database, err := db.NewDatabase(dbPath)
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

type exec struct {
	cmd string
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
