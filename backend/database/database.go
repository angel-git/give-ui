package database

import (
	"go.etcd.io/bbolt"
)

const configBucket = "config"
const dbName = "give-ui.config.db"

func CreateDatabase() *bbolt.DB {

	db, err := bbolt.Open(dbName, 0600, nil)
	if err != nil {
		panic("!!! Error opening database: " + err.Error())
	}
	createConfigBucket(db)
	return db
}

func GetValue(db *bbolt.DB, key string) string {
	var value = ""
	db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(configBucket))
		v := b.Get([]byte(key))
		if v != nil {
			value = string(v)
		}
		return nil
	})
	return value
}

func SaveValue(db *bbolt.DB, key string, value string) {
	db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(configBucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
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
		panic("!!! Error creating config bucket: " + err.Error())
	}
}
