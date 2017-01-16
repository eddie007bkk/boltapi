package main

import (
	"boltapi/db"
	"net/http"

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
		if database != nil {
			err := database.DB.Close()
			handleErr(err)
		}
		dbPath := dbsFolder + "/" + userReq.db
		database, err = db.NewDatabase(dbPath)
		handleErr(err)
		return

	}
}
