package db

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/sahandhnj/apiclient/pearl"
)

var db *bolt.DB
var BUCKET = []byte("pearl")

func Connect() (err error) {
	db, err = bolt.Open(".oyster/data.db", 0644, nil)

	return
}

func Close() (err error) {
	err = db.Close()

	return
}

func FindById(id string) (err error, value []byte) {
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BUCKET)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", BUCKET)
		}

		value = bucket.Get([]byte(id))
		fmt.Println(string(value))

		return nil
	})

	return
}

func Insert(pearl *pearl.Pearl) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BUCKET)
		if err != nil {
			return err
		}

		idB := []byte(pearl.ID)
		pearlB, err := json.Marshal(pearl)
		if err != nil {
			return err
		}

		c := bucket.Cursor()

		found := false
		for idB, v := c.First(); idB != nil; idB, v = c.Next() {
			if bytes.Equal(pearlB, v) {
				found = true
				break
			}
		}
		if found {
			return nil
		}

		return bucket.Put(idB, pearlB)
	})

	return
}
