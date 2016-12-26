package boltdb

import (
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// NewDatabase opens a new database
func NewDatabase(dbfile string) (d *Database, err error) {
	d = &Database{}
	d.DB, err = bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

// Database Struct
type Database struct {
	DB *bolt.DB
}

// Put inserts a key:value pair into the database
func (bt *Database) Put(bucket, key, val []byte) error {
	//dbPath := bt.db.Path()
	//log.Println("DB Info: ", reflect.TypeOf(dbPath), dbPath)
	err := bt.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		err = bucket.Put(key, val)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get retrieves a key:value pair from the database
func (bt *Database) Get(bucket, key []byte) (result []byte, err error) {
	bt.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			v := b.Get([]byte(key))
			if v != nil {
				result = make([]byte, len(v))
				copy(result, v)
			}
		} else {
			result = []byte("")
		}

		return nil
	})
	return
}

// CurrentDB retrieves the path of the current database
func (bt *Database) CurrentDB() string {
	dbPath := bt.DB.Path()
	dbName := strings.Split(dbPath, "/")
	return dbName[len(dbName)-1]
}
