package main

import (
	"boltapi/boltdb"
	"log"
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
	database.Get([]byte("bucket"), []byte("somekey"))
}
