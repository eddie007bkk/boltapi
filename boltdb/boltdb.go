package boltdb

import (
	"time"

	"github.com/boltdb/bolt"
)

// NewDatabase opens a new database
func NewDatabase(dbfile string) (d *Database, err error) {
	d = &Database{}
	d.db, err = bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

// Database Struct
type Database struct {
	db *bolt.DB
}

// Put inserts a key:value pair into the database
func (bt *Database) Put(bucket, key, val []byte) error {
	err := bt.db.Update(func(tx *bolt.Tx) error {
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
	bt.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		if v != nil {
			result = make([]byte, len(v))
			copy(result, v)
		}
		return nil
	})
	return
}
