package database

import (
	"internal/kvstore"
	"log"
)

type Database struct {
	dataDir string
	stores  map[string]*kvstore.KVStore
}

func (db *Database) Persist() {
	return
}

func (db *Database) GetOrCreateStore(name *string) *kvstore.KVStore {
	store, ok := db.stores[*name]
	if !ok {
		log.Printf("creating db %s\n", *name)
		store = kvstore.New(*name)
		db.stores[*name] = store
	} else {
		log.Printf("switched to db %s\n", *name)
	}
	return store
}

func New(dataDir *string) *Database {
	db := new(Database)
	db.dataDir = *dataDir
	db.stores = make(map[string]*kvstore.KVStore)
	log.Printf("instanciated new database at %s\n", *dataDir)
	return db
}
