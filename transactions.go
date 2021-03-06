package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

func txHandler(w http.ResponseWriter, r *http.Request) {
	var userReq userRequest
	userReq.GetRequest(r)

	//userReq.GetUserDB opens the requested databasse
	userReq.GetUserDB(r)

	if len(userReq.key) > 0 {
		switch r.Method {
		case "PUT":
			txPut(w, r, userReq)
		case "GET":
			txGet(w, r, userReq)
		case "DELETE":
			txDeleteKey(w, r, userReq)
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
		txDeleteDatabase(w, r, userReq)
	}
}

func txGet(w http.ResponseWriter, r *http.Request, userRequest userRequest) (err error) {
	if len(userRequest.key) > 0 {
		//User has specified a db, bucket & key
		res, err := database.Get([]byte(userRequest.bucket), []byte(userRequest.key))
		if err != nil {
			w.Write([]byte(`{"delete":` + err.Error() + `}`))
		} else {
			response := "{\"" + userRequest.key + "\":\"" + string(res) + "\"}"
			w.Write([]byte(response))
		}

	} else if len(userRequest.db) > 0 {
		//User has only specified a db, return data about this db
		res := database.BK.Stats()
		log.Println(reflect.TypeOf(res), res)

	} else {
		//User has not specified any db, return data about all DBs
	}
	return
}

func txPut(w http.ResponseWriter, r *http.Request, userRequest userRequest) (err error) {
	if len(userRequest.key) > 0 {
		//User has specified a db, bucket & key

		bucket := userRequest.bucket
		key := userRequest.key
		val, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(`{"delete":` + err.Error() + `}`))
		} else {
			err = database.Put([]byte(bucket), []byte(key), []byte(val))
			if err != nil {
				w.Write([]byte(`{"delete":` + err.Error() + `}`))
			} else {
				w.Write([]byte(`{"put":"ok"}`))
			}
		}

	} else if len(userRequest.db) > 0 {
		//User has only specified a db, open a new database
	} else {
		//User has not specified any db, do nothing
	}
	return
}

func txDeleteKey(w http.ResponseWriter, r *http.Request, userRequest userRequest) (err error) {
	if len(userRequest.key) > 0 {
		//User has specified a db, bucket & key
		bucket := userRequest.bucket
		key := userRequest.key
		err = database.DeleteKey([]byte(bucket), []byte(key))
		if err != nil {
			w.Write([]byte(`{"delete":` + err.Error() + `}`))
			return
		}
		w.Write([]byte(`{"delete":"ok")`))
		return
	}
	return
}

func txDeleteDatabase(w http.ResponseWriter, r *http.Request, userRequest userRequest) (err error) {
	if database != nil {
		err = database.DB.Close()

	}

	dbPath := dbsFolder + "/" + userRequest.db
	err = os.Remove(dbPath)
	if err != nil {
		w.Write([]byte(`{"delete":` + err.Error() + `}`))
		return
	}
	w.Write([]byte(`{"delete":"ok")`))
	return
}
