package database

import (
	"go.etcd.io/bbolt"
	"log"
)

const configBucket = "config"
const dbName = "give-ui.config.db"

func CreateDatabase() *bbolt.DB {
	db, err := bbolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	createConfigBucket(db)
	return db
}

func GetValue(db *bbolt.DB, key string) string {
	var value = ""
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(configBucket))
		v := b.Get([]byte(key))
		if v != nil {
			value = string(v)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading key [%s]: %s", key, err)
	}
	return value
}

func SaveValue(db *bbolt.DB, key string, value string) {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(configBucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		log.Fatalf("Error writing key [%s] with value [%s]: %s", key, value, err)
	}
}

func createConfigBucket(db *bbolt.DB) {
	err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(configBucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error creating the json bucket: %s", err)
	}
}
