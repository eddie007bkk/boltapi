package boltdb

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// New Database opens a new database
func NewDatabase(dbfile string) (d *database, err error) {
	d = &database{}
	d.db, err = bolt.Open(dbfile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	return
}

type database struct {
	db *bolt.DB
}

func (bt *database) Put(bucket, key, val []byte) error {
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

func (bt *database) Get(bucket, key []byte) {
	bt.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		log.Println("Val:", v)
		return nil
	})

}
