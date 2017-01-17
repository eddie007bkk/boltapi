package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	//Get variables {db}, {bucketName}, and {keyName} from user request URL
	var userRequest userRequest
	userRequest.GetRequest(r)

	reqURI := r.URL.RequestURI()
	reqURI = reqURI[1 : len(reqURI)-1]
	uri := strings.Split(reqURI, "/")

	userRequest.cmd = uri[len(uri)-1]
	switch userRequest.cmd {
	case "dbs":
		adminGetDatabases(w, r)
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
	return
}

func adminGetDatabases(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(dbsFolder)
	var dbList []string
	for _, f := range files {
		dbList = append(dbList, `"`+f.Name()+`"`)
	}
	res := `{"databases":[` + strings.Join(dbList, ",") + `]}`
	w.Write([]byte(res))
	return
}
